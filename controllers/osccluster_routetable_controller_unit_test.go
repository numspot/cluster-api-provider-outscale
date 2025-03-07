/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	infrastructurev1beta1 "github.com/outscale/cluster-api-provider-outscale/api/v1beta1"
	"github.com/outscale/cluster-api-provider-outscale/cloud/scope"
	"github.com/outscale/cluster-api-provider-outscale/cloud/services/security/mock_security"
	"github.com/outscale/cluster-api-provider-outscale/cloud/tag/mock_tag"
	osc "github.com/outscale/osc-sdk-go/v2"
	"github.com/stretchr/testify/require"
)

var (
	defaultRouteTableGatewayInitialize = infrastructurev1beta1.OscClusterSpec{
		Network: infrastructurev1beta1.OscNetwork{
			ClusterName: "test-cluster",
			Net: infrastructurev1beta1.OscNet{
				Name:    "test-net",
				IpRange: "10.0.0.0/16",
			},
			Subnets: []*infrastructurev1beta1.OscSubnet{
				{
					Name:          "test-subnet",
					IpSubnetRange: "10.0.0.0/24",
					SubregionName: "eu-west-2a",
				},
			},
			InternetService: infrastructurev1beta1.OscInternetService{
				Name: "test-internetservice",
			},
			RouteTables: []*infrastructurev1beta1.OscRouteTable{
				{
					Name: "test-routetable",
					Subnets: []string{
						"test-subnet",
					},
					Routes: []infrastructurev1beta1.OscRoute{
						{
							Name:        "test-route",
							TargetName:  "test-internetservice",
							TargetType:  "gateway",
							Destination: "0.0.0.0/0",
						},
					},
				},
			},
		},
	}
	defaultRouteTableGatewayReconcile = infrastructurev1beta1.OscClusterSpec{
		Network: infrastructurev1beta1.OscNetwork{
			ClusterName: "test-cluster",
			Net: infrastructurev1beta1.OscNet{
				Name:       "test-net",
				IpRange:    "10.0.0.0/16",
				ResourceId: "vpc-test-net-uid",
			},
			Subnets: []*infrastructurev1beta1.OscSubnet{
				{
					Name:          "test-subnet",
					IpSubnetRange: "10.0.0.0/24",
					SubregionName: "eu-west-2a",
					ResourceId:    "subnet-test-subnet-uid",
				},
			},
			InternetService: infrastructurev1beta1.OscInternetService{
				Name:       "test-internetservice",
				ResourceId: "igw-test-interneetservice-uid",
			},
			NatService: infrastructurev1beta1.OscNatService{
				Name:         "test-natservice",
				PublicIpName: "test-publicip",
				SubnetName:   "test-subnet",
				ResourceId:   "nat-test-natservice-uid",
			},
			RouteTables: []*infrastructurev1beta1.OscRouteTable{
				{
					Name: "test-routetable",
					Subnets: []string{
						"test-subnet",
					},
					ResourceId: "rtb-test-routetable-uid",
					Routes: []infrastructurev1beta1.OscRoute{
						{
							Name:        "test-route",
							TargetName:  "test-natservice",
							TargetType:  "nat",
							Destination: "0.0.0.0/0",
						},
					},
				},
			},
		},
	}

	defaultRouteTableNatInitialize = infrastructurev1beta1.OscClusterSpec{
		Network: infrastructurev1beta1.OscNetwork{
			ClusterName: "test-cluster",
			Net: infrastructurev1beta1.OscNet{
				Name:    "test-net",
				IpRange: "10.0.0.0/16",
			},
			Subnets: []*infrastructurev1beta1.OscSubnet{
				{
					Name:          "test-subnet",
					IpSubnetRange: "10.0.0.0/24",
					SubregionName: "eu-west-2a",
				},
			},
			InternetService: infrastructurev1beta1.OscInternetService{
				Name: "test-internetservice",
			},
			NatService: infrastructurev1beta1.OscNatService{
				Name:         "test-natservice",
				PublicIpName: "test-publicip",
				SubnetName:   "test-subnet",
			},
			RouteTables: []*infrastructurev1beta1.OscRouteTable{
				{
					Name: "test-routetable",
					Subnets: []string{
						"test-subnet",
					},
					Routes: []infrastructurev1beta1.OscRoute{
						{
							Name:        "test-route",
							TargetName:  "test-natservice",
							TargetType:  "nat",
							Destination: "0.0.0.0/0",
						},
					},
				},
			},
		},
	}

	defaultRouteTableNatReconcile = infrastructurev1beta1.OscClusterSpec{
		Network: infrastructurev1beta1.OscNetwork{
			ClusterName: "test-cluster",
			Net: infrastructurev1beta1.OscNet{
				Name:       "test-net",
				IpRange:    "10.0.0.0/16",
				ResourceId: "vpc-test-net-uid",
			},
			Subnets: []*infrastructurev1beta1.OscSubnet{
				{
					Name:          "test-subnet",
					IpSubnetRange: "10.0.0.0/24",
					SubregionName: "eu-west-2a",
					ResourceId:    "subnet-test-subnet-uid",
				},
			},
			InternetService: infrastructurev1beta1.OscInternetService{
				Name:       "test-internetservice",
				ResourceId: "igw-test-interneetservice-uid",
			},
			NatService: infrastructurev1beta1.OscNatService{
				Name:         "test-natservice",
				PublicIpName: "test-publicip",
				SubnetName:   "test-subnet",
				ResourceId:   "nat-test-natservice-uid",
			},
			RouteTables: []*infrastructurev1beta1.OscRouteTable{
				{
					Name: "test-routetable",
					Subnets: []string{
						"test-subnet",
					},
					ResourceId: "rtb-test-routetable-uid",
					Routes: []infrastructurev1beta1.OscRoute{
						{
							Name:        "test-route",
							TargetName:  "test-natservice",
							TargetType:  "nat",
							Destination: "0.0.0.0/0",
						},
					},
				},
			},
		},
	}

	defaultRouteTableGatewayNatInitialize = infrastructurev1beta1.OscClusterSpec{
		Network: infrastructurev1beta1.OscNetwork{
			ClusterName: "test-cluster",
			Net: infrastructurev1beta1.OscNet{
				Name:    "test-net",
				IpRange: "10.0.0.0/16",
			},
			Subnets: []*infrastructurev1beta1.OscSubnet{
				{
					Name:          "test-subnet",
					IpSubnetRange: "10.0.0.0/24",
					SubregionName: "eu-west-2a",
				},
			},
			InternetService: infrastructurev1beta1.OscInternetService{
				Name: "test-internetservice",
			},
			NatService: infrastructurev1beta1.OscNatService{
				Name:         "test-natservice",
				PublicIpName: "test-publicip",
				SubnetName:   "test-subnet",
			},
			RouteTables: []*infrastructurev1beta1.OscRouteTable{
				{
					Name: "test-routetable",
					Subnets: []string{
						"test-subnet",
					},
					Routes: []infrastructurev1beta1.OscRoute{
						{
							Name:        "test-route-nat",
							TargetName:  "test-natservice",
							TargetType:  "nat",
							Destination: "0.0.0.0/0",
						},
						{
							Name:        "test-route-igw",
							TargetName:  "test-internetservice",
							TargetType:  "gateway",
							Destination: "0.0.0.0/0",
						},
					},
				},
			},
		},
	}

	defaultRouteTableGatewayNatReconcile = infrastructurev1beta1.OscClusterSpec{
		Network: infrastructurev1beta1.OscNetwork{
			ClusterName: "test-cluster",
			Net: infrastructurev1beta1.OscNet{
				Name:       "test-net",
				IpRange:    "10.0.0.0/16",
				ResourceId: "vpc-test-net",
			},
			Subnets: []*infrastructurev1beta1.OscSubnet{
				{
					Name:          "test-subnet",
					IpSubnetRange: "10.0.0.0/24",
					SubregionName: "eu-west-2a",
					ResourceId:    "subnet-test-subnet-uid",
				},
			},
			InternetService: infrastructurev1beta1.OscInternetService{
				Name:       "test-internetservice",
				ResourceId: "igw-test-internetservice-uid",
			},
			NatService: infrastructurev1beta1.OscNatService{
				Name:         "test-natservice",
				PublicIpName: "test-publicip",
				SubnetName:   "test-subnet",
				ResourceId:   "nat-test-natservice-uid",
			},
			RouteTables: []*infrastructurev1beta1.OscRouteTable{
				{
					Name: "test-routetable",
					Subnets: []string{
						"test-subnet",
					},
					ResourceId: "rtb-test-routetable-uid",
					Routes: []infrastructurev1beta1.OscRoute{
						{
							Name:        "test-route-nat",
							TargetName:  "test-natservice",
							TargetType:  "nat",
							Destination: "0.0.0.0/0",
						},
						{
							Name:        "test-route-igw",
							TargetName:  "test-internetservice",
							TargetType:  "gateway",
							Destination: "0.0.0.0/0",
						},
					},
				},
			},
		},
	}
)

// SetupWithRouteTableMock set routeTableMock with clusterScope and osccluster
func SetupWithRouteTableMock(t *testing.T, name string, spec infrastructurev1beta1.OscClusterSpec) (clusterScope *scope.ClusterScope, ctx context.Context, mockOscRouteTableInterface *mock_security.MockOscRouteTableInterface, mockOscTagInterface *mock_tag.MockOscTagInterface) {
	clusterScope = Setup(t, name, spec)
	mockCtrl := gomock.NewController(t)
	mockOscRouteTableInterface = mock_security.NewMockOscRouteTableInterface(mockCtrl)
	mockOscTagInterface = mock_tag.NewMockOscTagInterface(mockCtrl)
	ctx = context.Background()
	return clusterScope, ctx, mockOscRouteTableInterface, mockOscTagInterface
}

// TestGettRouteTableResourceId has several tests to cover the code of the function getRouteTableResourceId
func TestGetRouteTableResourceId(t *testing.T) {
	routeTableTestCases := []struct {
		name                          string
		spec                          infrastructurev1beta1.OscClusterSpec
		expRouteTablesFound           bool
		expGetRouteTableResourceIdErr error
	}{
		{
			name:                          "get RouteTableId",
			spec:                          defaultRouteTableGatewayInitialize,
			expRouteTablesFound:           true,
			expGetRouteTableResourceIdErr: nil,
		},
		{
			name:                          "can not get RouteTableId",
			spec:                          defaultRouteTableGatewayInitialize,
			expRouteTablesFound:           false,
			expGetRouteTableResourceIdErr: errors.New("test-routetable-uid does not exist"),
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope := Setup(t, rttc.name, rttc.spec)
			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				if rttc.expRouteTablesFound {
					routeTablesRef.ResourceMap[routeTableName] = routeTableId
				}
				routeTableResourceId, err := getRouteTableResourceId(routeTableName, clusterScope)
				if rttc.expGetRouteTableResourceIdErr != nil {
					require.EqualError(t, err, rttc.expGetRouteTableResourceIdErr.Error(), "GetRouteTableResourceId() should return the same error")
				} else {
					require.NoError(t, err)
				}
				t.Logf("Find routeTableResourceId %s\n", routeTableResourceId)
			}
		})
	}
}

// TestGettRouteResourceId has several tests to cover the code of the function getRouteResourceId
func TestGetRouteResourceId(t *testing.T) {
	routeTestCases := []struct {
		name                     string
		spec                     infrastructurev1beta1.OscClusterSpec
		expRouteFound            bool
		expGetRouteResourceIdErr error
	}{
		{
			name:                     "get RouteId",
			spec:                     defaultRouteTableGatewayInitialize,
			expRouteFound:            true,
			expGetRouteResourceIdErr: nil,
		},
		{
			name:                     "can not get RouteId",
			spec:                     defaultRouteTableGatewayInitialize,
			expRouteFound:            false,
			expGetRouteResourceIdErr: errors.New("test-route-uid does not exist"),
		},
	}
	for _, rtc := range routeTestCases {
		t.Run(rtc.name, func(t *testing.T) {
			clusterScope := Setup(t, rtc.name, rtc.spec)
			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)
			routeTablesSpec := rtc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routesSpec := routeTableSpec.Routes
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				for _, routeSpec := range routesSpec {
					routeName := routeSpec.Name + "-uid"
					if rtc.expRouteFound {
						routeRef.ResourceMap[routeName] = routeTableId
					}
					routeResourceId, err := getRouteResourceId(routeName, clusterScope)
					if rtc.expGetRouteResourceIdErr != nil {
						require.EqualError(t, err, rtc.expGetRouteResourceIdErr.Error(), "GetRouteResourceId() should return the same error")
					} else {
						require.NoError(t, err)
					}
					t.Logf("Find routeResourceId %s\n", routeResourceId)
				}
			}
		})
	}
}

// TestCheckRouteTableSubnetOscAssociateResourceName has several tests to cover the code of the func checkRouteTableSubnetOscAssociateResourceName
func TestCheckRouteTableSubnetOscAssociateResourceName(t *testing.T) {
	routeTableTestCases := []struct {
		name                                                string
		spec                                                infrastructurev1beta1.OscClusterSpec
		expCheckRouteTableSubnetOscAssociateResourceNameErr error
	}{
		{
			name: "check work without net, routetable and route spec (with default values)",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{},
			},
			expCheckRouteTableSubnetOscAssociateResourceNameErr: nil,
		},
		{
			name: "check routetable association with subnet",
			spec: defaultRouteTableGatewayInitialize,
			expCheckRouteTableSubnetOscAssociateResourceNameErr: nil,
		},
		{
			name: "check routetable association with bad subnet",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{
							Name: "test-routetable",
							Subnets: []string{
								"test-subnet-test",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "0.0.0.0/0",
								},
							},
						},
					},
				},
			},
			expCheckRouteTableSubnetOscAssociateResourceNameErr: errors.New("subnet test-subnet-test-uid does not exist in routeTable"),
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope := Setup(t, rttc.name, rttc.spec)
			err := checkRouteTableSubnetOscAssociateResourceName(clusterScope)
			if rttc.expCheckRouteTableSubnetOscAssociateResourceNameErr != nil {
				require.EqualError(t, err, rttc.expCheckRouteTableSubnetOscAssociateResourceNameErr.Error(), "CheckRouteTableSubnetOscAssociateResourceName() should return the same error")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestCheckRouteTableFormatParameters has several tests to cover the code of the func checkRouteTableFormatParameters
func TestCheckRouteTableFormatParameters(t *testing.T) {
	routeTableTestCases := []struct {
		name                                  string
		spec                                  infrastructurev1beta1.OscClusterSpec
		expCheckRouteTableFormatParametersErr error
	}{
		{
			name: "check work without net, routable and route spec (with default values)",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{},
			},
			expCheckRouteTableFormatParametersErr: nil,
		},
		{
			name:                                  "check routetable format",
			spec:                                  defaultRouteTableGatewayInitialize,
			expCheckRouteTableFormatParametersErr: nil,
		},
		{
			name: "check Bad Name routetable",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{
							Name: "test-routetable@test",
							Subnets: []string{
								"test-subnet",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "0.0.0.0/0",
								},
							},
						},
					},
				},
			},
			expCheckRouteTableFormatParametersErr: errors.New("Invalid Tag Name"),
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope := Setup(t, rttc.name, rttc.spec)
			_, err := checkRouteTableFormatParameters(clusterScope)
			if rttc.expCheckRouteTableFormatParametersErr != nil {
				require.EqualError(t, err, rttc.expCheckRouteTableFormatParametersErr.Error(), "CheckRouteTableFormatParameters() should return the same error")
			} else {
				require.NoError(t, err)
			}
			t.Logf("find all routetablename ")
		})
	}
}

// TestCheckRouteFormatParameters has several tests to cover the code of the func checkRouteFormatParameters
func TestCheckRouteFormatParameters(t *testing.T) {
	routeTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expCheckRouteFormatParametersErr error
	}{
		{
			name: "check work without net, routetable and route spec (with default values)",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{},
			},
			expCheckRouteFormatParametersErr: nil,
		},
		{
			name:                             "check route format",
			spec:                             defaultRouteTableGatewayInitialize,
			expCheckRouteFormatParametersErr: nil,
		},
		{
			name: "check Bad Name route",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{
							Name: "test-routetable",
							Subnets: []string{
								"test-subnet",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route@test",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "0.0.0.0/0",
								},
							},
						},
					},
				},
			},
			expCheckRouteFormatParametersErr: errors.New("Invalid Tag Name"),
		},
		{
			name: "check Bad Ip Range IP route",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{
							Name: "test-routetable",
							Subnets: []string{
								"test-subnet",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "10.0.0.256/16",
								},
							},
						},
					},
				},
			},
			expCheckRouteFormatParametersErr: errors.New("invalid CIDR address: 10.0.0.256/16"),
		},
		{
			name: "check Bad Ip Range IP route",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{
							Name: "test-routetable",
							Subnets: []string{
								"test-subnet",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "10.0.0.0/36",
								},
							},
						},
					},
				},
			},
			expCheckRouteFormatParametersErr: errors.New("invalid CIDR address: 10.0.0.0/36"),
		},
	}
	for _, rtc := range routeTestCases {
		t.Run(rtc.name, func(t *testing.T) {
			clusterScope := Setup(t, rtc.name, rtc.spec)

			_, err := checkRouteFormatParameters(clusterScope)
			if rtc.expCheckRouteFormatParametersErr != nil {
				require.EqualError(t, err, rtc.expCheckRouteFormatParametersErr.Error(), "CheckRouteFormatParameters() should return the same error")
			} else {
				require.NoError(t, err)
			}
			t.Logf("find all routeName")
		})
	}
}

// TestCheckRouteTableOscDuplicateName has several tests to cover the code of the func checkRouteTableOscDuplicateName
func TestCheckRouteTableOscDuplicateName(t *testing.T) {
	routeTableTestCases := []struct {
		name                                  string
		spec                                  infrastructurev1beta1.OscClusterSpec
		expCheckRouteTableOscDuplicateNameErr error
	}{
		{
			name:                                  "get no duplicate routeTable Name",
			spec:                                  defaultRouteTableGatewayInitialize,
			expCheckRouteTableOscDuplicateNameErr: nil,
		},
		{
			name: "get duplicate routeTable Name",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{
							Name: "test-routetable",
							Subnets: []string{
								"test-subnet",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "0.0.0.0/0",
								},
							},
						},
						{
							Name: "test-routetable",
							Subnets: []string{
								"test-subnet",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "0.0.0.0/0",
								},
							},
						},
					},
				},
			},
			expCheckRouteTableOscDuplicateNameErr: errors.New("test-routetable already exist"),
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope := Setup(t, rttc.name, rttc.spec)
			err := checkRouteTableOscDuplicateName(clusterScope)
			if rttc.expCheckRouteTableOscDuplicateNameErr != nil {
				require.EqualError(t, err, rttc.expCheckRouteTableOscDuplicateNameErr.Error(), "checkRouteTableOscDuplicateName() should return the same error")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestCheckRouteOscDuplicateName has several tests to cover the code of the func checkRouteOscDuplicateName
func TestCheckRouteOscDuplicateName(t *testing.T) {
	routeTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expCheckRouteOscDuplicateNameErr error
	}{
		{
			name: "check work without net, routetable and route spec (with default values)",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{},
			},
			expCheckRouteOscDuplicateNameErr: nil,
		},
		{
			name:                             "check route duplicate name",
			spec:                             defaultRouteTableGatewayInitialize,
			expCheckRouteOscDuplicateNameErr: nil,
		},
		{
			name:                             "get no route duplicate name",
			spec:                             defaultRouteTableGatewayInitialize,
			expCheckRouteOscDuplicateNameErr: nil,
		},
		{
			name: "check route duplicate  internet service name",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{
							Name: "test-routetable",
							Subnets: []string{
								"test-subnet",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "0.0.0.0/0",
								},
								{
									Name:        "test-route",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "0.0.0.0/0",
								},
							},
						},
					},
				},
			},
			expCheckRouteOscDuplicateNameErr: errors.New("test-route already exist"),
		},
		{
			name: "check route duplicate  nat service name",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{
							Name: "test-routetable",
							Subnets: []string{
								"test-subnet",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route",
									TargetName:  "test-natservice",
									TargetType:  "nat",
									Destination: "0.0.0.0/0",
								},
								{
									Name:        "test-route",
									TargetName:  "test-natservice",
									TargetType:  "nat",
									Destination: "0.0.0.0/0",
								},
							},
						},
					},
				},
			},
			expCheckRouteOscDuplicateNameErr: errors.New("test-route already exist"),
		},
		{
			name: "check no routetable",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{},
					},
				},
			},
			expCheckRouteOscDuplicateNameErr: nil,
		},
	}
	for _, rtc := range routeTestCases {
		t.Run(rtc.name, func(t *testing.T) {
			clusterScope := Setup(t, rtc.name, rtc.spec)
			err := checkRouteOscDuplicateName(clusterScope)
			if rtc.expCheckRouteOscDuplicateNameErr != nil {
				require.EqualError(t, err, rtc.expCheckRouteOscDuplicateNameErr.Error(), "CheckRouteOscDuplicateName() should return the same error")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestReconcilerRouteCreate has several tests to cover the code of the function reconcilerRouteCreate
func TestReconcilerRouteCreate(t *testing.T) {
	routeTestCases := []struct {
		name                         string
		spec                         infrastructurev1beta1.OscClusterSpec
		expRouteFound                bool
		expTagFound                  bool
		expInternetServiceFound      bool
		expNatServiceFound           bool
		expCreateRouteFound          bool
		expCreateRouteErr            error
		expGetRouteTableFromRouteErr error
		expReadTagErr                error
		expReconcileRouteErr         error
	}{
		{
			name:                         "create route with internet service (first time reconcile loop)",
			spec:                         defaultRouteTableGatewayInitialize,
			expRouteFound:                false,
			expInternetServiceFound:      true,
			expNatServiceFound:           false,
			expCreateRouteFound:          true,
			expTagFound:                  false,
			expCreateRouteErr:            nil,
			expGetRouteTableFromRouteErr: nil,
			expReadTagErr:                nil,
			expReconcileRouteErr:         nil,
		},
		{
			name:                         "create route with natservice (first time reconcile loop)",
			spec:                         defaultRouteTableNatInitialize,
			expRouteFound:                false,
			expInternetServiceFound:      false,
			expNatServiceFound:           true,
			expTagFound:                  false,
			expCreateRouteFound:          true,
			expCreateRouteErr:            nil,
			expGetRouteTableFromRouteErr: nil,
			expReadTagErr:                nil,
			expReconcileRouteErr:         nil,
		},
		{
			name:                         "create multi route  (first time reconcile loop)",
			spec:                         defaultRouteTableGatewayNatInitialize,
			expRouteFound:                false,
			expTagFound:                  false,
			expInternetServiceFound:      true,
			expNatServiceFound:           true,
			expCreateRouteFound:          true,
			expCreateRouteErr:            nil,
			expGetRouteTableFromRouteErr: nil,
			expReadTagErr:                nil,
			expReconcileRouteErr:         nil,
		},
		{
			name:                         "user delete route without cluster-api",
			spec:                         defaultRouteTableNatReconcile,
			expRouteFound:                false,
			expTagFound:                  false,
			expInternetServiceFound:      false,
			expNatServiceFound:           true,
			expCreateRouteFound:          true,
			expCreateRouteErr:            nil,
			expGetRouteTableFromRouteErr: nil,
			expReadTagErr:                nil,
			expReconcileRouteErr:         nil,
		},
		{
			name:                         "failed to create route",
			spec:                         defaultRouteTableNatInitialize,
			expRouteFound:                false,
			expTagFound:                  false,
			expInternetServiceFound:      false,
			expNatServiceFound:           true,
			expCreateRouteFound:          false,
			expCreateRouteErr:            errors.New("CreateRoute generic error"),
			expGetRouteTableFromRouteErr: nil,
			expReadTagErr:                nil,
			expReconcileRouteErr:         errors.New("cannot create route: CreateRoute generic error"),
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, mockOscTagInterface := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)

			var associateRouteTableId string
			var resourceId string
			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				tag := osc.Tag{
					ResourceId: &routeTableId,
				}
				if rttc.expTagFound {
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(&tag, rttc.expReadTagErr)
				}
				routeTablesRef.ResourceMap[routeTableName] = routeTableId
				associateRouteTableId = routeTableId
				routesSpec := routeTableSpec.Routes
				for _, routeSpec := range routesSpec {
					destinationIpRange := routeSpec.Destination
					resourceType := routeSpec.TargetType
					if resourceType == "gateway" {
						resourceId = internetServiceId
					} else {
						resourceId = natServiceId
					}

					route := osc.CreateRouteResponse{
						RouteTable: &osc.RouteTable{
							RouteTableId: &routeTableId,
						},
					}

					readRouteTables := osc.ReadRouteTablesResponse{
						RouteTables: &[]osc.RouteTable{
							*route.RouteTable,
						},
					}
					readRouteTable := *readRouteTables.RouteTables
					if rttc.expRouteFound {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(&readRouteTable[0], rttc.expGetRouteTableFromRouteErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(nil, rttc.expGetRouteTableFromRouteErr)
					}
					if rttc.expCreateRouteFound {
						mockOscRouteTableInterface.
							EXPECT().
							CreateRoute(gomock.Any(), gomock.Eq(destinationIpRange), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(route.RouteTable, rttc.expCreateRouteErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							CreateRoute(gomock.Any(), gomock.Eq(destinationIpRange), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(nil, rttc.expCreateRouteErr)
					}
					reconcileRoute, err := reconcileRoute(ctx, clusterScope, routeSpec, routeTableName, mockOscRouteTableInterface)
					if rttc.expReconcileRouteErr != nil {
						require.EqualError(t, err, rttc.expReconcileRouteErr.Error(), "reconcileRoute() should return the same error")
					} else {
						require.NoError(t, err)
					}
					t.Logf("find reconcileRoute %v\n", reconcileRoute)
				}
			}
		})
	}
}

// TestReconcileRouteGet has several tests to cover the code of the function reconcileRouteGet
func TestReconcileRouteGet(t *testing.T) {
	routeTestCases := []struct {
		name                         string
		spec                         infrastructurev1beta1.OscClusterSpec
		expRouteFound                bool
		expTagFound                  bool
		expInternetServiceFound      bool
		expNatServiceFound           bool
		expGetRouteTableFromRouteErr error
		expReadTagErr                error
		expReconcileRouteErr         error
	}{
		{
			name:                         "check reconcile multi route (second time reconcile loop)",
			spec:                         defaultRouteTableGatewayNatReconcile,
			expRouteFound:                true,
			expTagFound:                  false,
			expInternetServiceFound:      true,
			expNatServiceFound:           true,
			expGetRouteTableFromRouteErr: nil,
			expReadTagErr:                nil,
			expReconcileRouteErr:         nil,
		},
		{
			name:                         "check reconcile route with natservice (second time reconcile loop)",
			spec:                         defaultRouteTableNatReconcile,
			expRouteFound:                true,
			expTagFound:                  false,
			expInternetServiceFound:      false,
			expNatServiceFound:           true,
			expGetRouteTableFromRouteErr: nil,
			expReadTagErr:                nil,
			expReconcileRouteErr:         nil,
		},
		{
			name:                         "failed to get route",
			spec:                         defaultRouteTableNatInitialize,
			expRouteFound:                false,
			expTagFound:                  false,
			expInternetServiceFound:      false,
			expNatServiceFound:           true,
			expGetRouteTableFromRouteErr: errors.New("GetRouteTableFromRoute generic error"),
			expReadTagErr:                nil,
			expReconcileRouteErr:         errors.New("cannot get route table: GetRouteTableFromRoute generic error"),
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, mockOscTagInterface := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)
			var associateRouteTableId string
			var resourceId string

			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				routeTablesRef.ResourceMap[routeTableName] = routeTableId
				tag := osc.Tag{
					ResourceId: &routeTableId,
				}
				if rttc.expTagFound {
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(&tag, rttc.expReadTagErr)
				}
				associateRouteTableId = routeTableId
				routesSpec := routeTableSpec.Routes
				for _, routeSpec := range routesSpec {
					resourceType := routeSpec.TargetType
					if resourceType == "gateway" {
						resourceId = internetServiceId
					} else {
						resourceId = natServiceId
					}

					route := osc.CreateRouteResponse{
						RouteTable: &osc.RouteTable{
							RouteTableId: &routeTableId,
						},
					}

					readRouteTables := osc.ReadRouteTablesResponse{
						RouteTables: &[]osc.RouteTable{
							*route.RouteTable,
						},
					}
					readRouteTable := *readRouteTables.RouteTables
					if rttc.expRouteFound {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(&readRouteTable[0], rttc.expGetRouteTableFromRouteErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(nil, rttc.expGetRouteTableFromRouteErr)
					}
					reconcileRoute, err := reconcileRoute(ctx, clusterScope, routeSpec, routeTableName, mockOscRouteTableInterface)
					if rttc.expReconcileRouteErr != nil {
						require.EqualError(t, err, rttc.expReconcileRouteErr.Error(), "reconcileRoute() should return the same error")
					} else {
						require.NoError(t, err)
					}
					t.Logf("find reconcileRoute %v\n", reconcileRoute)
				}
			}
		})
	}
}

// TestReconcileRouteResourceId has several tests to cover the code of the function reconcileRouteResourceId
func TestReconcileRouteResourceId(t *testing.T) {
	routeTestCases := []struct {
		name                    string
		spec                    infrastructurev1beta1.OscClusterSpec
		expInternetServiceFound bool
		expNatServiceFound      bool
		expTagFound             bool
		expReadTagErr           error
		expReconcileRouteErr    error
	}{
		{
			name:                    "natService does not exist",
			spec:                    defaultRouteTableNatInitialize,
			expInternetServiceFound: false,
			expNatServiceFound:      false,
			expTagFound:             false,
			expReadTagErr:           nil,
			expReconcileRouteErr:    errors.New("test-natservice-uid does not exist"),
		},
		{
			name:                    "internetService does not exist",
			spec:                    defaultRouteTableGatewayInitialize,
			expInternetServiceFound: false,
			expNatServiceFound:      false,
			expTagFound:             false,
			expReadTagErr:           nil,
			expReconcileRouteErr:    errors.New("test-internetservice-uid does not exist"),
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, mockOscTagInterface := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				tag := osc.Tag{
					ResourceId: &routeTableId,
				}
				if rttc.expTagFound {
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(&tag, rttc.expReadTagErr)
				}
				routesSpec := routeTableSpec.Routes
				for _, routeSpec := range routesSpec {
					reconcileRoute, err := reconcileRoute(ctx, clusterScope, routeSpec, routeTableName, mockOscRouteTableInterface)
					if rttc.expReconcileRouteErr != nil {
						require.EqualError(t, err, rttc.expReconcileRouteErr.Error(), "reconcileRoute() should return the same error")
					} else {
						require.NoError(t, err)
					}
					t.Logf("find reconcileRoute %v\n", reconcileRoute)
				}
			}
		})
	}
}

// TestReconcileRouteTableCreate has several tests to cover the code of the function reconcileRouteTableCreate
func TestReconcileRouteTableCreate(t *testing.T) {
	routeTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expNetFound                      bool
		expSubnetFound                   bool
		expRouteFound                    bool
		expRouteTableFound               bool
		expInternetServiceFound          bool
		expNatServiceFound               bool
		expCreateRouteFound              bool
		expCreateRouteTableFound         bool
		expLinkRouteTableFound           bool
		expTagFound                      bool
		expCreateRouteErr                error
		expCreateRouteTableErr           error
		expLinkRouteTableErr             error
		expGetRouteTableFromRouteErr     error
		expGetRouteTableIdsFromNetIdsErr error
		expReadTagErr                    error
		expReconcileRouteTableErr        error
	}{
		{
			name:                             "create routetable with internet service route (first time reconcile loop)",
			spec:                             defaultRouteTableGatewayInitialize,
			expNetFound:                      true,
			expSubnetFound:                   true,
			expRouteFound:                    false,
			expRouteTableFound:               false,
			expInternetServiceFound:          true,
			expNatServiceFound:               false,
			expCreateRouteFound:              true,
			expCreateRouteTableFound:         true,
			expLinkRouteTableFound:           true,
			expTagFound:                      false,
			expCreateRouteErr:                nil,
			expCreateRouteTableErr:           nil,
			expLinkRouteTableErr:             nil,
			expGetRouteTableFromRouteErr:     nil,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReadTagErr:                    nil,
			expReconcileRouteTableErr:        nil,
		},
		{
			name:                             "failed to create route",
			spec:                             defaultRouteTableGatewayInitialize,
			expNetFound:                      true,
			expSubnetFound:                   true,
			expRouteFound:                    false,
			expRouteTableFound:               false,
			expTagFound:                      false,
			expInternetServiceFound:          true,
			expNatServiceFound:               false,
			expCreateRouteFound:              true,
			expCreateRouteTableFound:         true,
			expLinkRouteTableFound:           true,
			expCreateRouteErr:                errors.New("CreateRoute generic error"),
			expCreateRouteTableErr:           nil,
			expLinkRouteTableErr:             nil,
			expGetRouteTableFromRouteErr:     nil,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReadTagErr:                    nil,
			expReconcileRouteTableErr:        errors.New("cannot create route: CreateRoute generic error"),
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, mockOscTagInterface := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := rttc.spec.Network.Net.Name + "-uid"
			netId := "vpc-" + netName
			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			if rttc.expNetFound {
				netRef.ResourceMap[netName] = netId
			}
			clusterName := rttc.spec.Network.ClusterName + "-uid"
			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			linkRouteTableRef := clusterScope.GetLinkRouteTablesRef()
			if len(linkRouteTableRef) == 0 {
				linkRouteTableRef = make(map[string][]string)
			}
			subnetRef := clusterScope.GetSubnetRef()
			subnetRef.ResourceMap = make(map[string]string)

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)

			var associateRouteTableId string
			var routeTableIds []string
			var resourceId string

			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				tag := osc.Tag{
					ResourceId: &routeTableId,
				}
				if rttc.expTagFound {
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(&tag, rttc.expReadTagErr)
				} else {
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(nil, rttc.expReadTagErr)
				}
				routeTableIds = append(routeTableIds, routeTableId)
				linkRouteTableId := "eipalloc-" + routeTableName
				subnetsSpec := routeTableSpec.Subnets
				for _, subnet := range subnetsSpec {
					subnetName := subnet + "-uid"
					subnetId := "subnet-" + subnetName

					if rttc.expSubnetFound {
						subnetRef.ResourceMap[subnetName] = subnetId
					}

					if rttc.expLinkRouteTableFound {
						linkRouteTableRef[routeTableName] = []string{linkRouteTableId}
					}

					routeTable := osc.CreateRouteTableResponse{
						RouteTable: &osc.RouteTable{
							RouteTableId: &routeTableId,
						},
					}

					linkRouteTable := osc.LinkRouteTableResponse{
						LinkRouteTableId: &linkRouteTableId,
					}

					readRouteTables := osc.ReadRouteTablesResponse{
						RouteTables: &[]osc.RouteTable{
							*routeTable.RouteTable,
						},
					}
					readRouteTable := *readRouteTables.RouteTables
					if rttc.expRouteTableFound {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
							Return(routeTableIds, rttc.expGetRouteTableIdsFromNetIdsErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
							Return(nil, rttc.expGetRouteTableIdsFromNetIdsErr)
					}
					if rttc.expCreateRouteTableFound {
						associateRouteTableId = routeTableId
						routeTablesRef.ResourceMap[routeTableName] = routeTableId
						mockOscRouteTableInterface.
							EXPECT().
							CreateRouteTable(gomock.Any(), gomock.Eq(netId), gomock.Eq(clusterName), gomock.Eq(routeTableName)).
							Return(routeTable.RouteTable, rttc.expCreateRouteTableErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							CreateRouteTable(gomock.Any(), gomock.Eq(netId), gomock.Eq(netName), gomock.Eq(routeTableName)).
							Return(nil, rttc.expCreateRouteTableErr)
					}

					if rttc.expLinkRouteTableFound {
						mockOscRouteTableInterface.
							EXPECT().
							LinkRouteTable(gomock.Any(), gomock.Eq(routeTableId), gomock.Eq(subnetId)).
							Return(*linkRouteTable.LinkRouteTableId, rttc.expLinkRouteTableErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							LinkRouteTable(gomock.Any(), gomock.Eq(routeTableId), gomock.Eq(subnetId)).
							Return("", rttc.expLinkRouteTableErr)
					}

					routesSpec := routeTableSpec.Routes
					for _, routeSpec := range routesSpec {
						destinationIpRange := routeSpec.Destination
						resourceType := routeSpec.TargetType
						if resourceType == "gateway" {
							resourceId = internetServiceId
						} else {
							resourceId = natServiceId
						}

						route := osc.CreateRouteResponse{
							RouteTable: &osc.RouteTable{
								RouteTableId: &routeTableId,
							},
						}
						if rttc.expRouteFound {
							mockOscRouteTableInterface.
								EXPECT().
								GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
								Return(&readRouteTable[0], rttc.expGetRouteTableFromRouteErr)
						} else {
							mockOscRouteTableInterface.
								EXPECT().
								GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
								Return(nil, rttc.expGetRouteTableFromRouteErr)
						}
						if rttc.expCreateRouteFound {
							mockOscRouteTableInterface.
								EXPECT().
								CreateRoute(gomock.Any(), gomock.Eq(destinationIpRange), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
								Return(route.RouteTable, rttc.expCreateRouteErr)
						} else {
							mockOscRouteTableInterface.
								EXPECT().
								CreateRoute(gomock.Any(), gomock.Eq(destinationIpRange), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
								Return(nil, rttc.expCreateRouteErr)
						}
					}
					reconcileRouteTable, err := reconcileRouteTable(ctx, clusterScope, mockOscRouteTableInterface, mockOscTagInterface)
					if rttc.expReconcileRouteTableErr != nil {
						require.EqualError(t, err, rttc.expReconcileRouteTableErr.Error(), "reconcileRouteTable() should return the same error")
					} else {
						require.NoError(t, err)
					}
					t.Logf("find reconcileRoute %v\n", reconcileRouteTable)
				}
			}
		})
	}
}

// reconcileRouteTableGet has several tests to cover the code of the function reconcileRouteTableGet
func TestReconcileRouteTableGet(t *testing.T) {
	routeTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expNetFound                      bool
		expTagFound                      bool
		expSubnetFound                   bool
		expRouteTableFound               bool
		expInternetServiceFound          bool
		expNatServiceFound               bool
		expGetRouteTableIdsFromNetIdsErr error
		expReadTagErr                    error
		expReconcileRouteTableErr        error
	}{
		{
			name:                             "check reconcile routetable  with internet service route (second time reconcile loop)",
			spec:                             defaultRouteTableGatewayReconcile,
			expNetFound:                      true,
			expSubnetFound:                   true,
			expRouteTableFound:               true,
			expInternetServiceFound:          true,
			expNatServiceFound:               false,
			expTagFound:                      true,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReadTagErr:                    nil,
			expReconcileRouteTableErr:        nil,
		},
		{
			name:                             "failed to get routetable",
			spec:                             defaultRouteTableGatewayInitialize,
			expNetFound:                      true,
			expSubnetFound:                   true,
			expRouteTableFound:               false,
			expInternetServiceFound:          true,
			expNatServiceFound:               false,
			expTagFound:                      false,
			expGetRouteTableIdsFromNetIdsErr: errors.New("GetRouteTableIdsFromNetIds generic errors"),
			expReadTagErr:                    nil,
			expReconcileRouteTableErr:        errors.New("list route tables: GetRouteTableIdsFromNetIds generic errors"),
		},
		{
			name:                             "create routetable with natservice (first time reconcile loop)",
			spec:                             defaultRouteTableNatInitialize,
			expNetFound:                      true,
			expSubnetFound:                   true,
			expTagFound:                      true,
			expRouteTableFound:               true,
			expInternetServiceFound:          false,
			expNatServiceFound:               false,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReadTagErr:                    nil,
			expReconcileRouteTableErr:        nil,
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, mockOscTagInterface := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := rttc.spec.Network.Net.Name + "-uid"
			netId := "vpc-" + netName
			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			if rttc.expNetFound {
				netRef.ResourceMap[netName] = netId
			}

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			subnetRef := clusterScope.GetSubnetRef()
			subnetRef.ResourceMap = make(map[string]string)

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			var routeTableIds []string

			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				routeTableIds = append(routeTableIds, routeTableId)
				subnetsSpec := routeTableSpec.Subnets
				for _, subnet := range subnetsSpec {
					subnetName := subnet + "-uid"
					subnetId := "subnet-" + subnetName
					tag := osc.Tag{
						ResourceId: &subnetId,
					}
					if rttc.expTagFound {
						if rttc.expRouteTableFound {
							mockOscTagInterface.
								EXPECT().
								ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
								Return(&tag, rttc.expReadTagErr)
						}
					}
					if rttc.expSubnetFound {
						subnetRef.ResourceMap[subnetName] = subnetId
					}
					if rttc.expRouteTableFound {
						routeTablesRef.ResourceMap[routeTableName] = routeTableId
					}

					if rttc.expRouteTableFound {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
							Return(routeTableIds, rttc.expGetRouteTableIdsFromNetIdsErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
							Return(nil, rttc.expGetRouteTableIdsFromNetIdsErr)
					}
				}

				reconcileRouteTable, err := reconcileRouteTable(ctx, clusterScope, mockOscRouteTableInterface, mockOscTagInterface)
				if rttc.expReconcileRouteTableErr != nil {
					require.EqualError(t, err, rttc.expReconcileRouteTableErr.Error(), "reconcileRouteTable() should return the same error")
				} else {
					require.NoError(t, err)
				}
				t.Logf("find reconcileRoute %v\n", reconcileRouteTable)
			}
		})
	}
}

// TestReconcileRouteTableResourceId has several tests to cover the code of the function reconcileRouteTable
func TestReconcileRouteTableResourceId(t *testing.T) {
	routeTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expTagFound                      bool
		expNetFound                      bool
		expReadTagErr                    error
		expReconcileRouteTableErr        error
		expGetRouteTableIdsFromNetIdsErr error
	}{
		{
			name:                             "net does not exist",
			spec:                             defaultRouteTableGatewayInitialize,
			expTagFound:                      false,
			expNetFound:                      false,
			expReadTagErr:                    nil,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReconcileRouteTableErr:        errors.New("test-net-uid does not exist"),
		},
		{
			name:                             "failed to get tag",
			spec:                             defaultRouteTableGatewayInitialize,
			expTagFound:                      true,
			expNetFound:                      true,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReadTagErr:                    errors.New("ReadTag generic error"),
			expReconcileRouteTableErr:        errors.New("cannot get tag: ReadTag generic error"),
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, mockOscTagInterface := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			netName := rttc.spec.Network.Net.Name + "-uid"
			netId := "vpc-" + netName
			if rttc.expNetFound {
				netRef.ResourceMap[netName] = netId
			}
			routeTablesSpec := rttc.spec.Network.RouteTables
			var routeTableIds []string
			if rttc.expTagFound {
				for _, routeTableSpec := range routeTablesSpec {
					routeTableName := routeTableSpec.Name + "-uid"
					routeTableId := "rtb-" + routeTableName
					routeTableIds = append(routeTableIds, routeTableId)
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(nil, rttc.expReadTagErr)
				}

				mockOscRouteTableInterface.
					EXPECT().
					GetRouteTableIdsFromNetIds(gomock.Any(), netId).
					Return(routeTableIds, rttc.expGetRouteTableIdsFromNetIdsErr)
			}
			reconcileRouteTable, err := reconcileRouteTable(ctx, clusterScope, mockOscRouteTableInterface, mockOscTagInterface)
			if rttc.expReconcileRouteTableErr != nil {
				require.EqualError(t, err, rttc.expReconcileRouteTableErr.Error(), "reconcileRouteTable() should return the same error")
			} else {
				require.NoError(t, err)
			}
			t.Logf("find reconcileRoute %v\n", reconcileRouteTable)
		})
	}
}

// TestReconcileCreateRouteTable has several tests to cover the code of the function reconcileRouteTable
func TestReconcileCreateRouteTable(t *testing.T) {
	routeTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expTagFound                      bool
		expCreateRouteTableErr           error
		expGetRouteTableIdsFromNetIdsErr error
		expReadTagErr                    error
		expReconcileRouteTableErr        error
	}{
		{
			name: "failed to create routeTable",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{
					Net: infrastructurev1beta1.OscNet{
						Name:    "test-net",
						IpRange: "10.0.0.0/16",
					},
					Subnets: []*infrastructurev1beta1.OscSubnet{
						{
							Name:          "test-subnet",
							IpSubnetRange: "10.0.0.0/24",
							SubregionName: "eu-west-2a",
						},
					},
					InternetService: infrastructurev1beta1.OscInternetService{
						Name: "test-internetservice",
					},
					RouteTables: []*infrastructurev1beta1.OscRouteTable{
						{
							Name: "test-routetable",
							Subnets: []string{
								"test-subnet",
							},
							Routes: []infrastructurev1beta1.OscRoute{
								{
									Name:        "test-route",
									TargetName:  "test-internetservice",
									TargetType:  "gateway",
									Destination: "0.0.0.0/0",
								},
							},
						},
					},
				},
			},
			expCreateRouteTableErr:           errors.New("CreateRouteTable generic error"),
			expGetRouteTableIdsFromNetIdsErr: nil,
			expTagFound:                      false,
			expReadTagErr:                    nil,
			expReconcileRouteTableErr:        errors.New("cannot create routetable: CreateRouteTable generic error"),
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, mockOscTagInterface := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := rttc.spec.Network.Net.Name + "-uid"
			clusterName := rttc.spec.Network.ClusterName + "-uid"
			netId := "vpc-" + netName
			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			netRef.ResourceMap[netName] = netId

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			subnetRef := clusterScope.GetSubnetRef()
			subnetRef.ResourceMap = make(map[string]string)
			var routeTableIds []string

			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				routeTableIds = append(routeTableIds, routeTableId)
				tag := osc.Tag{
					ResourceId: &routeTableId,
				}
				if rttc.expTagFound {
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(&tag, rttc.expReadTagErr)
				} else {
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(nil, rttc.expReadTagErr)
				}

				subnetsSpec := routeTableSpec.Subnets
				for _, subnet := range subnetsSpec {
					subnetName := subnet + "-uid"
					subnetId := "subnet-" + subnetName
					subnetRef.ResourceMap[subnetName] = subnetId
					mockOscRouteTableInterface.
						EXPECT().
						GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
						Return(routeTableIds, rttc.expGetRouteTableIdsFromNetIdsErr)

					mockOscRouteTableInterface.
						EXPECT().
						CreateRouteTable(gomock.Any(), gomock.Eq(netId), gomock.Eq(clusterName), gomock.Eq(routeTableName)).
						Return(nil, rttc.expCreateRouteTableErr)
				}
				reconcileRouteTable, err := reconcileRouteTable(ctx, clusterScope, mockOscRouteTableInterface, mockOscTagInterface)
				if rttc.expReconcileRouteTableErr != nil {
					require.EqualError(t, err, rttc.expReconcileRouteTableErr.Error(), "reconcileRouteTable() should return the same error")
				} else {
					require.NoError(t, err)
				}
				t.Logf("find reconcileRoute %v\n", reconcileRouteTable)
			}
		})
	}
}

// TestReconcileRouteTableLink has several tests to cover the code of the function reconcileRouteTable
func TestReconcileRouteTableLink(t *testing.T) {
	routeTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expSubnetFound                   bool
		expTagFound                      bool
		expLinkRouteTableFound           bool
		expCreateRouteTableErr           error
		expLinkRouteTableErr             error
		expGetRouteTableIdsFromNetIdsErr error
		expReadTagErr                    error
		expReconcileRouteTableErr        error
	}{
		{
			name:                             "failed to link routeTable",
			spec:                             defaultRouteTableGatewayInitialize,
			expSubnetFound:                   true,
			expTagFound:                      false,
			expLinkRouteTableFound:           true,
			expCreateRouteTableErr:           nil,
			expLinkRouteTableErr:             errors.New("LinkRouteTable generic error"),
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReadTagErr:                    nil,
			expReconcileRouteTableErr:        errors.New("cannot link routetable with net: LinkRouteTable generic error"),
		},
		{
			name:                             "failed to get subnet",
			spec:                             defaultRouteTableGatewayInitialize,
			expSubnetFound:                   false,
			expLinkRouteTableFound:           false,
			expTagFound:                      false,
			expCreateRouteTableErr:           nil,
			expLinkRouteTableErr:             nil,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReadTagErr:                    nil,
			expReconcileRouteTableErr:        errors.New("test-subnet-uid does not exist"),
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, mockOscTagInterface := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := rttc.spec.Network.Net.Name + "-uid"
			netId := "vpc-" + netName
			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			netRef.ResourceMap[netName] = netId

			clusterName := rttc.spec.Network.ClusterName + "-uid"
			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			subnetRef := clusterScope.GetSubnetRef()
			subnetRef.ResourceMap = make(map[string]string)

			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				tag := osc.Tag{
					ResourceId: &routeTableId,
				}
				if rttc.expTagFound {
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(&tag, rttc.expReadTagErr)
				} else {
					mockOscTagInterface.
						EXPECT().
						ReadTag(gomock.Any(), gomock.Eq("Name"), gomock.Eq(routeTableName)).
						Return(nil, rttc.expReadTagErr)
				}
				subnetsSpec := routeTableSpec.Subnets
				for _, subnet := range subnetsSpec {
					subnetName := subnet + "-uid"
					subnetId := "subnet-" + subnetName
					if rttc.expSubnetFound {
						subnetRef.ResourceMap[subnetName] = subnetId
					}

					routeTable := osc.CreateRouteTableResponse{
						RouteTable: &osc.RouteTable{
							RouteTableId: &routeTableId,
						},
					}

					mockOscRouteTableInterface.
						EXPECT().
						CreateRouteTable(gomock.Any(), gomock.Eq(netId), gomock.Eq(clusterName), gomock.Eq(routeTableName)).
						Return(routeTable.RouteTable, rttc.expCreateRouteTableErr)
					if rttc.expLinkRouteTableFound {
						mockOscRouteTableInterface.
							EXPECT().
							LinkRouteTable(gomock.Any(), gomock.Eq(routeTableId), gomock.Eq(subnetId)).
							Return("", rttc.expLinkRouteTableErr)
					}
					mockOscRouteTableInterface.
						EXPECT().
						GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
						Return(nil, rttc.expGetRouteTableIdsFromNetIdsErr)
				}
				reconcileRouteTable, err := reconcileRouteTable(ctx, clusterScope, mockOscRouteTableInterface, mockOscTagInterface)
				if rttc.expReconcileRouteTableErr != nil {
					require.EqualError(t, err, rttc.expReconcileRouteTableErr.Error(), "reconcileRouteTable() should return the same error")
				} else {
					require.NoError(t, err)
				}
				t.Logf("find reconcileRoute %v\n", reconcileRouteTable)
			}
		})
	}
}

// TestReconcileDeleteRouteDelete has several tests to cover the code of the function reconcileDeleteRoute
func TestReconcileDeleteRouteDelete(t *testing.T) {
	routeTestCases := []struct {
		name                         string
		spec                         infrastructurev1beta1.OscClusterSpec
		expRouteFound                bool
		expInternetServiceFound      bool
		expNatServiceFound           bool
		expDeleteRouteErr            error
		expGetRouteTableFromRouteErr error
		expReconcileDeleteRouteErr   error
	}{
		{
			name:                         "delete Route with internetservice (first time reconcile loop)",
			spec:                         defaultRouteTableGatewayInitialize,
			expRouteFound:                true,
			expInternetServiceFound:      true,
			expNatServiceFound:           false,
			expDeleteRouteErr:            nil,
			expGetRouteTableFromRouteErr: nil,
			expReconcileDeleteRouteErr:   nil,
		},
		{
			name:                         "delete Route with natservice (first time reconcile loop)",
			spec:                         defaultRouteTableNatReconcile,
			expRouteFound:                true,
			expInternetServiceFound:      false,
			expNatServiceFound:           true,
			expDeleteRouteErr:            nil,
			expGetRouteTableFromRouteErr: nil,
			expReconcileDeleteRouteErr:   nil,
		},
		{
			name:                         "delete Route with internetservice  and gatewayservice (first time reconcile loop)",
			spec:                         defaultRouteTableGatewayNatReconcile,
			expRouteFound:                true,
			expInternetServiceFound:      true,
			expNatServiceFound:           true,
			expDeleteRouteErr:            nil,
			expGetRouteTableFromRouteErr: nil,
			expReconcileDeleteRouteErr:   nil,
		},
		{
			name:                         "failed to delete route",
			spec:                         defaultRouteTableGatewayInitialize,
			expRouteFound:                true,
			expInternetServiceFound:      true,
			expNatServiceFound:           false,
			expDeleteRouteErr:            errors.New("DeleteRoute generic error"),
			expGetRouteTableFromRouteErr: nil,
			expReconcileDeleteRouteErr:   errors.New("cannot delete route: DeleteRoute generic error"),
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, _ := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			var associateRouteTableId string
			var resourceId string
			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				routeTablesRef.ResourceMap[routeTableName] = routeTableId
				associateRouteTableId = routeTableId
				routesSpec := routeTableSpec.Routes
				for _, routeSpec := range routesSpec {
					destinationIpRange := routeSpec.Destination
					resourceType := routeSpec.TargetType
					routeName := routeSpec.Name + "-uid"
					routeRef.ResourceMap[routeName] = routeTableId
					if resourceType == "gateway" {
						resourceId = internetServiceId
					} else {
						resourceId = natServiceId
					}
					route := osc.CreateRouteResponse{
						RouteTable: &osc.RouteTable{
							RouteTableId: &routeTableId,
						},
					}

					readRouteTables := osc.ReadRouteTablesResponse{
						RouteTables: &[]osc.RouteTable{
							*route.RouteTable,
						},
					}

					readRouteTable := *readRouteTables.RouteTables
					if rttc.expRouteFound {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(&readRouteTable[0], rttc.expGetRouteTableFromRouteErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(nil, rttc.expGetRouteTableFromRouteErr)
					}
					mockOscRouteTableInterface.
						EXPECT().
						DeleteRoute(gomock.Any(), gomock.Eq(destinationIpRange), gomock.Eq(routeTableId)).
						Return(rttc.expDeleteRouteErr)

					reconcileDeleteRoute, err := reconcileDeleteRoute(ctx, clusterScope, routeSpec, routeTableName, mockOscRouteTableInterface)
					if rttc.expReconcileDeleteRouteErr != nil {
						require.EqualError(t, err, rttc.expReconcileDeleteRouteErr.Error(), "reconcileDeleteRoute() should return the same error")
					} else {
						require.NoError(t, err)
					}
					t.Logf("Find reconcileDeleteRoute %v\n", reconcileDeleteRoute)
				}
			}
		})
	}
}

// TestReconcileDeleteRouteGet has several tests to cover the code of the function reconcileDeleteRoute
func TestReconcileDeleteRouteGet(t *testing.T) {
	routeTestCases := []struct {
		name                         string
		spec                         infrastructurev1beta1.OscClusterSpec
		expRouteFound                bool
		expInternetServiceFound      bool
		expNatServiceFound           bool
		expGetRouteTableFromRouteErr error
		expReconcileDeleteRouteErr   error
	}{
		{
			name:                         "failed to get route",
			spec:                         defaultRouteTableGatewayInitialize,
			expRouteFound:                false,
			expInternetServiceFound:      true,
			expNatServiceFound:           false,
			expGetRouteTableFromRouteErr: errors.New("GetRouteTable generic error"),
			expReconcileDeleteRouteErr:   errors.New("checking route table: GetRouteTable generic error"),
		},
		{
			name:                         "remove finalizer (user delete route without cluster-api)",
			spec:                         defaultRouteTableGatewayInitialize,
			expRouteFound:                false,
			expInternetServiceFound:      true,
			expNatServiceFound:           true,
			expGetRouteTableFromRouteErr: nil,
			expReconcileDeleteRouteErr:   nil,
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, _ := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			var associateRouteTableId string
			var resourceId string
			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				routeTablesRef.ResourceMap[routeTableName] = routeTableId
				associateRouteTableId = routeTableId
				routesSpec := routeTableSpec.Routes
				for _, routeSpec := range routesSpec {
					resourceType := routeSpec.TargetType
					routeName := routeSpec.Name + "-uid"
					routeRef.ResourceMap[routeName] = routeTableId

					if resourceType == "gateway" {
						resourceId = internetServiceId
					} else {
						resourceId = natServiceId
					}
					route := osc.CreateRouteResponse{
						RouteTable: &osc.RouteTable{
							RouteTableId: &routeTableId,
						},
					}

					readRouteTables := osc.ReadRouteTablesResponse{
						RouteTables: &[]osc.RouteTable{
							*route.RouteTable,
						},
					}

					readRouteTable := *readRouteTables.RouteTables
					if rttc.expRouteFound {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(&readRouteTable[0], rttc.expGetRouteTableFromRouteErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(nil, rttc.expGetRouteTableFromRouteErr)
					}

					reconcileDeleteRoute, err := reconcileDeleteRoute(ctx, clusterScope, routeSpec, routeTableName, mockOscRouteTableInterface)
					if rttc.expReconcileDeleteRouteErr != nil {
						require.EqualError(t, err, rttc.expReconcileDeleteRouteErr.Error(), "reconcileDeleteRoute() should return the same error")
					} else {
						require.NoError(t, err)
					}
					t.Logf("Find reconcileDeleteRoute %v\n", reconcileDeleteRoute)
				}
			}
		})
	}
}

// TestReconcileDeleteRouteResourceId has several tests to cover the code of the function reconcileDeleteRoute
func TestReconcileDeleteRouteResourceId(t *testing.T) {
	routeTestCases := []struct {
		name                       string
		spec                       infrastructurev1beta1.OscClusterSpec
		expInternetServiceFound    bool
		expNatServiceFound         bool
		expReconcileDeleteRouteErr error
	}{
		{
			name:                       "natService does not exist",
			spec:                       defaultRouteTableNatReconcile,
			expInternetServiceFound:    false,
			expNatServiceFound:         false,
			expReconcileDeleteRouteErr: errors.New("test-natservice-uid does not exist"),
		},
		{
			name:                       "internetService does not exist",
			spec:                       defaultRouteTableGatewayInitialize,
			expInternetServiceFound:    false,
			expNatServiceFound:         false,
			expReconcileDeleteRouteErr: errors.New("test-internetservice-uid does not exist"),
		},
		{
			name:                       "route does not exist",
			spec:                       defaultRouteTableGatewayInitialize,
			expInternetServiceFound:    true,
			expNatServiceFound:         true,
			expReconcileDeleteRouteErr: errors.New("test-route-uid does not exist"),
		},
	}
	for _, rttc := range routeTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, _ := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			internetServiceId := "igw-" + internetServiceName
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routesSpec := routeTableSpec.Routes
				for _, routeSpec := range routesSpec {
					reconcileDeleteRoute, err := reconcileDeleteRoute(ctx, clusterScope, routeSpec, routeTableName, mockOscRouteTableInterface)
					if rttc.expReconcileDeleteRouteErr != nil {
						require.EqualError(t, err, rttc.expReconcileDeleteRouteErr.Error(), "reconcileDeleteRoute() should return the same error")
					} else {
						require.NoError(t, err)
					}
					t.Logf("Find reconcileDeleteRoute %v\n", reconcileDeleteRoute)
				}
			}
		})
	}
}

// TestReconcileDeleteRouteTableDelete has several tests to cover the code of the function reconcileDeleteRouteTable
func TestReconcileDeleteRouteTableDeleteWithoutSpec(t *testing.T) {
	routeTableTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expUnlinkRouteTableErr           error
		expDeleteRouteErr                error
		expDeleteRouteTableErr           error
		expGetRouteTableFromRouteErr     error
		expGetRouteTableIdsFromNetIdsErr error
		expReconcileDeleteRouteTableErr  error
	}{
		{
			name:                             "delete Routetable with internetservice route (first time reconcile loop) without spec (with default values)",
			expUnlinkRouteTableErr:           nil,
			expDeleteRouteErr:                nil,
			expDeleteRouteTableErr:           nil,
			expGetRouteTableFromRouteErr:     nil,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReconcileDeleteRouteTableErr:  nil,
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, _ := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := "cluster-api-net-uid"
			netId := "vpc-" + netName
			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			netRef.ResourceMap[netName] = netId

			internetServiceName := "cluster-api-internetservice-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			internetServiceRef.ResourceMap[internetServiceName] = internetServiceId

			natServiceName := "cluster-api-natservice-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			natServiceRef.ResourceMap[natServiceName] = natServiceId

			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			linkRouteTableRef := clusterScope.GetLinkRouteTablesRef()
			if len(linkRouteTableRef) == 0 {
				linkRouteTableRef = make(map[string][]string)
			}
			var associateRouteTableId string
			var resourceId string
			var routeTableIds []string
			routeTableName := "cluster-api-routetable-kw-uid"
			routeTableId := "rtb-" + routeTableName
			routeTableIds = append(routeTableIds, routeTableId)
			linkRouteTableId := "eipalloc-" + routeTableName
			routeTablesRef.ResourceMap[routeTableName] = routeTableId
			associateRouteTableId = routeTableId

			linkRouteTableRef[routeTableName] = []string{linkRouteTableId}
			clusterScope.SetLinkRouteTablesRef(linkRouteTableRef)
			mockOscRouteTableInterface.
				EXPECT().
				GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
				Return(routeTableIds, rttc.expGetRouteTableIdsFromNetIdsErr)
			mockOscRouteTableInterface.
				EXPECT().
				UnlinkRouteTable(gomock.Any(), gomock.Eq(linkRouteTableId)).
				Return(rttc.expUnlinkRouteTableErr)

			mockOscRouteTableInterface.
				EXPECT().
				DeleteRouteTable(gomock.Any(), gomock.Eq(routeTableId)).
				Return(rttc.expDeleteRouteTableErr)

			destinationIpRange := "0.0.0.0/0"
			resourceType := "nat"
			routeName := "cluster-api-route-kw-uid"
			routeRef.ResourceMap[routeName] = routeTableId

			if resourceType == "gateway" {
				resourceId = internetServiceId
			} else {
				resourceId = natServiceId
			}
			route := osc.CreateRouteResponse{
				RouteTable: &osc.RouteTable{
					RouteTableId: &routeTableId,
				},
			}

			readRouteTables := osc.ReadRouteTablesResponse{
				RouteTables: &[]osc.RouteTable{
					*route.RouteTable,
				},
			}

			readRouteTable := *readRouteTables.RouteTables
			mockOscRouteTableInterface.
				EXPECT().
				GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
				Return(&readRouteTable[0], rttc.expGetRouteTableFromRouteErr)

			mockOscRouteTableInterface.
				EXPECT().
				DeleteRoute(gomock.Any(), gomock.Eq(destinationIpRange), gomock.Eq(routeTableId)).
				Return(rttc.expDeleteRouteErr)
			reconcileDeleteRouteTable, err := reconcileDeleteRouteTable(ctx, clusterScope, mockOscRouteTableInterface)
			if rttc.expReconcileDeleteRouteTableErr != nil {
				require.EqualError(t, err, rttc.expReconcileDeleteRouteTableErr.Error(), "reconcileDeleteRouteTable() should return the same error")
			} else {
				require.NoError(t, err)
			}
			t.Logf("Find reconcileDeleteRouteTable %v\n", reconcileDeleteRouteTable)
		})
	}
}

// TestReconcileDeleteRouteTableDelete has several tests to cover the code of the function reconcileDeleteRouteTable
func TestReconcileDeleteRouteTableDelete(t *testing.T) {
	routeTableTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expNetFound                      bool
		expRouteFound                    bool
		expRouteTableFound               bool
		expInternetServiceFound          bool
		expNatServiceFound               bool
		expLinkRouteTableFound           bool
		expUnlinkRouteTableErr           error
		expDeleteRouteErr                error
		expDeleteRouteTableErr           error
		expGetRouteTableFromRouteErr     error
		expGetRouteTableIdsFromNetIdsErr error
		expReconcileDeleteRouteTableErr  error
	}{
		{
			name:                             "delete Routetable with internetservice route (first time reconcile loop)",
			spec:                             defaultRouteTableGatewayInitialize,
			expNetFound:                      true,
			expRouteFound:                    true,
			expRouteTableFound:               true,
			expInternetServiceFound:          true,
			expNatServiceFound:               false,
			expLinkRouteTableFound:           true,
			expUnlinkRouteTableErr:           nil,
			expDeleteRouteErr:                nil,
			expDeleteRouteTableErr:           nil,
			expGetRouteTableFromRouteErr:     nil,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReconcileDeleteRouteTableErr:  nil,
		},
		{
			name:                             "failed to delete routetable",
			spec:                             defaultRouteTableGatewayInitialize,
			expNetFound:                      true,
			expRouteFound:                    true,
			expRouteTableFound:               true,
			expInternetServiceFound:          true,
			expNatServiceFound:               false,
			expLinkRouteTableFound:           true,
			expUnlinkRouteTableErr:           nil,
			expDeleteRouteErr:                nil,
			expDeleteRouteTableErr:           errors.New("DeleteRoutetable generic error"),
			expGetRouteTableFromRouteErr:     nil,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReconcileDeleteRouteTableErr:  errors.New("cannot delete routeTable: DeleteRoutetable generic error"),
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, _ := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := rttc.spec.Network.Net.Name + "-uid"
			netId := "vpc-" + netName
			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			if rttc.expNetFound {
				netRef.ResourceMap[netName] = netId
			}

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)

			linkRouteTableRef := clusterScope.GetLinkRouteTablesRef()
			if len(linkRouteTableRef) == 0 {
				linkRouteTableRef = make(map[string][]string)
			}

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			var associateRouteTableId string
			var resourceId string
			var routeTableIds []string
			routeTablesSpec := rttc.spec.Network.RouteTables

			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				routeTableIds = append(routeTableIds, routeTableId)
				linkRouteTableId := "eipalloc-" + routeTableName
				if rttc.expRouteTableFound {
					routeTablesRef.ResourceMap[routeTableName] = routeTableId
					associateRouteTableId = routeTableId
				}

				if rttc.expLinkRouteTableFound {
					linkRouteTableRef[routeTableName] = []string{linkRouteTableId}
					clusterScope.SetLinkRouteTablesRef(linkRouteTableRef)
				}

				if rttc.expRouteTableFound {
					mockOscRouteTableInterface.
						EXPECT().
						GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
						Return(routeTableIds, rttc.expGetRouteTableIdsFromNetIdsErr)
				} else {
					mockOscRouteTableInterface.
						EXPECT().
						GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
						Return(nil, rttc.expGetRouteTableIdsFromNetIdsErr)
				}
				mockOscRouteTableInterface.
					EXPECT().
					UnlinkRouteTable(gomock.Any(), gomock.Eq(linkRouteTableId)).
					Return(rttc.expUnlinkRouteTableErr)
				mockOscRouteTableInterface.
					EXPECT().
					DeleteRouteTable(gomock.Any(), gomock.Eq(routeTableId)).
					Return(rttc.expDeleteRouteTableErr)

				routesSpec := routeTableSpec.Routes
				for _, routeSpec := range routesSpec {
					destinationIpRange := routeSpec.Destination
					resourceType := routeSpec.TargetType
					routeName := routeSpec.Name + "-uid"
					if rttc.expRouteFound {
						routeRef.ResourceMap[routeName] = routeTableId
					}
					if resourceType == "gateway" {
						resourceId = internetServiceId
					} else {
						resourceId = natServiceId
					}
					route := osc.CreateRouteResponse{
						RouteTable: &osc.RouteTable{
							RouteTableId: &routeTableId,
						},
					}

					readRouteTables := osc.ReadRouteTablesResponse{
						RouteTables: &[]osc.RouteTable{
							*route.RouteTable,
						},
					}

					readRouteTable := *readRouteTables.RouteTables
					if rttc.expRouteTableFound {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(&readRouteTable[0], rttc.expGetRouteTableFromRouteErr)
					} else {
						mockOscRouteTableInterface.
							EXPECT().
							GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
							Return(nil, rttc.expGetRouteTableFromRouteErr)
					}
					mockOscRouteTableInterface.
						EXPECT().
						DeleteRoute(gomock.Any(), gomock.Eq(destinationIpRange), gomock.Eq(routeTableId)).
						Return(rttc.expDeleteRouteErr)
				}
			}
			reconcileDeleteRouteTable, err := reconcileDeleteRouteTable(ctx, clusterScope, mockOscRouteTableInterface)
			if rttc.expReconcileDeleteRouteTableErr != nil {
				require.EqualError(t, err, rttc.expReconcileDeleteRouteTableErr.Error(), "reconcileDeleteRouteTable() should return the same error")
			} else {
				require.NoError(t, err)
			}
			t.Logf("Find reconcileDeleteRouteTable %v\n", reconcileDeleteRouteTable)
		})
	}
}

// TestReconcileDeleteRouteTableGet has several tests to cover the code of the function reconcileDeleteRouteTable
func TestReconcileDeleteRouteTableGet(t *testing.T) {
	routeTableTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expNetFound                      bool
		expRouteFound                    bool
		expRouteTableFound               bool
		expInternetServiceFound          bool
		expNatServiceFound               bool
		expGetRouteTableIdsFromNetIdsErr error
		expReconcileDeleteRouteTableErr  error
	}{
		{
			name:                             "failed to get routetable",
			spec:                             defaultRouteTableGatewayInitialize,
			expNetFound:                      true,
			expRouteTableFound:               false,
			expInternetServiceFound:          false,
			expNatServiceFound:               false,
			expGetRouteTableIdsFromNetIdsErr: errors.New("GetRouteTableIdsFromNetIds generic error"),
			expReconcileDeleteRouteTableErr:  errors.New("GetRouteTableIdsFromNetIds generic error"),
		},
		{
			name:                             "remove finalizer (delete routetable without cluster-api)",
			spec:                             defaultRouteTableGatewayInitialize,
			expNetFound:                      true,
			expRouteTableFound:               false,
			expInternetServiceFound:          false,
			expNatServiceFound:               false,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReconcileDeleteRouteTableErr:  nil,
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, _ := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := rttc.spec.Network.Net.Name + "-uid"
			netId := "vpc-" + netName
			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			if rttc.expNetFound {
				netRef.ResourceMap[netName] = netId
			}

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			if rttc.expInternetServiceFound {
				internetServiceRef.ResourceMap[internetServiceName] = internetServiceId
			}

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			if rttc.expNatServiceFound {
				natServiceRef.ResourceMap[natServiceName] = natServiceId
			}

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				if rttc.expRouteTableFound {
					routeTablesRef.ResourceMap[routeTableName] = routeTableId
				}

				if rttc.expRouteTableFound {
					mockOscRouteTableInterface.
						EXPECT().
						GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
						Return([]string{routeTableId}, rttc.expGetRouteTableIdsFromNetIdsErr)
				} else {
					mockOscRouteTableInterface.
						EXPECT().
						GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
						Return(nil, rttc.expGetRouteTableIdsFromNetIdsErr)
				}
			}
			reconcileDeleteRouteTable, err := reconcileDeleteRouteTable(ctx, clusterScope, mockOscRouteTableInterface)
			if rttc.expReconcileDeleteRouteTableErr != nil {
				require.EqualError(t, err, rttc.expReconcileDeleteRouteTableErr.Error(), "reconcileDeleteRouteTable() should return the same error")
			} else {
				require.NoError(t, err)
			}
			t.Logf("Find reconcileDeleteRouteTable %v\n", reconcileDeleteRouteTable)
		})
	}
}

// TestReconcileDeleteRouteTableUnlink has several tests to cover the code of the function reconcileDeleteRouteTable
func TestReconcileDeleteRouteTableUnlink(t *testing.T) {
	routeTableTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expUnlinkRouteTableErr           error
		expDeleteRouteErr                error
		expGetRouteTableFromRouteErr     error
		expGetRouteTableIdsFromNetIdsErr error
		expReconcileDeleteRouteTableErr  error
	}{
		{
			name:                             "failed to unlink routetable",
			spec:                             defaultRouteTableGatewayInitialize,
			expUnlinkRouteTableErr:           errors.New("UnlinkRouteTable generic error"),
			expDeleteRouteErr:                nil,
			expGetRouteTableFromRouteErr:     nil,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReconcileDeleteRouteTableErr:  errors.New("cannot unlink routeTable: UnlinkRouteTable generic error"),
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, _ := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := rttc.spec.Network.Net.Name + "-uid"
			netId := "vpc-" + netName
			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			netRef.ResourceMap[netName] = netId

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			internetServiceRef.ResourceMap[internetServiceName] = internetServiceId

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)
			natServiceRef.ResourceMap[natServiceName] = natServiceId

			linkRouteTableRef := make(map[string][]string)

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)
			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)

			var associateRouteTableId string
			var resourceId string
			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				linkRouteTableId := "eipalloc-" + routeTableName
				routeTablesRef.ResourceMap[routeTableName] = routeTableId
				associateRouteTableId = routeTableId

				linkRouteTableRef[routeTableName] = []string{linkRouteTableId}
				clusterScope.SetLinkRouteTablesRef(linkRouteTableRef)
				mockOscRouteTableInterface.
					EXPECT().
					GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
					Return([]string{routeTableId}, rttc.expGetRouteTableIdsFromNetIdsErr)

				mockOscRouteTableInterface.
					EXPECT().
					UnlinkRouteTable(gomock.Any(), gomock.Eq(linkRouteTableId)).
					Return(rttc.expUnlinkRouteTableErr)

				routesSpec := routeTableSpec.Routes
				for _, routeSpec := range routesSpec {
					destinationIpRange := routeSpec.Destination
					resourceType := routeSpec.TargetType
					routeName := routeSpec.Name + "-uid"
					routeRef.ResourceMap[routeName] = routeTableId
					if resourceType == "gateway" {
						resourceId = internetServiceId
					} else {
						resourceId = natServiceId
					}
					route := osc.CreateRouteResponse{
						RouteTable: &osc.RouteTable{
							RouteTableId: &routeTableId,
						},
					}

					readRouteTables := osc.ReadRouteTablesResponse{
						RouteTables: &[]osc.RouteTable{
							*route.RouteTable,
						},
					}

					readRouteTable := *readRouteTables.RouteTables
					mockOscRouteTableInterface.
						EXPECT().
						GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
						Return(&readRouteTable[0], rttc.expGetRouteTableFromRouteErr)
					mockOscRouteTableInterface.
						EXPECT().
						DeleteRoute(gomock.Any(), gomock.Eq(destinationIpRange), gomock.Eq(routeTableId)).
						Return(rttc.expDeleteRouteErr).AnyTimes()
				}
			}
			reconcileDeleteRouteTable, err := reconcileDeleteRouteTable(ctx, clusterScope, mockOscRouteTableInterface)
			if rttc.expReconcileDeleteRouteTableErr != nil {
				require.EqualError(t, err, rttc.expReconcileDeleteRouteTableErr.Error(), "reconcileDeleteRouteTable() should return the same error")
			} else {
				require.NoError(t, err)
			}
			t.Logf("Find reconcileDeleteRouteTable %v\n", reconcileDeleteRouteTable)
		})
	}
}

// TestReconcileDeleteRouteDeleteRouteTable has several tests to cover the code of the function reconcileDeleteRouteTable
func TestReconcileDeleteRouteDeleteRouteTable(t *testing.T) {
	routeTableTestCases := []struct {
		name                             string
		spec                             infrastructurev1beta1.OscClusterSpec
		expDeleteRouteErr                error
		expGetRouteTableFromRouteErr     error
		expGetRouteTableIdsFromNetIdsErr error
		expReconcileDeleteRouteTableErr  error
	}{
		{
			name:                             "failed to delete route",
			spec:                             defaultRouteTableGatewayInitialize,
			expDeleteRouteErr:                errors.New("DeleteRoute generic error"),
			expGetRouteTableFromRouteErr:     nil,
			expGetRouteTableIdsFromNetIdsErr: nil,
			expReconcileDeleteRouteTableErr:  errors.New("cannot delete route: DeleteRoute generic error"),
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, _ := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := rttc.spec.Network.Net.Name + "-uid"
			netId := "vpc-" + netName
			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			netRef.ResourceMap[netName] = netId

			internetServiceName := rttc.spec.Network.InternetService.Name + "-uid"
			internetServiceId := "igw-" + internetServiceName
			internetServiceRef := clusterScope.GetInternetServiceRef()
			internetServiceRef.ResourceMap = make(map[string]string)
			internetServiceRef.ResourceMap[internetServiceName] = internetServiceId

			natServiceName := rttc.spec.Network.NatService.Name + "-uid"
			natServiceId := "nat-" + natServiceName
			natServiceRef := clusterScope.GetNatServiceRef()
			natServiceRef.ResourceMap = make(map[string]string)

			routeRef := clusterScope.GetRouteRef()
			routeRef.ResourceMap = make(map[string]string)

			linkRouteTableRef := make(map[string][]string)

			routeTablesRef := clusterScope.GetRouteTablesRef()
			routeTablesRef.ResourceMap = make(map[string]string)

			var associateRouteTableId string
			var resourceId string
			routeTablesSpec := rttc.spec.Network.RouteTables
			for _, routeTableSpec := range routeTablesSpec {
				routeTableName := routeTableSpec.Name + "-uid"
				routeTableId := "rtb-" + routeTableName
				linkRouteTableId := "eipalloc-" + routeTableName
				routeTablesRef.ResourceMap[routeTableName] = routeTableId
				associateRouteTableId = routeTableId

				linkRouteTableRef[routeTableName] = []string{linkRouteTableId}
				clusterScope.SetLinkRouteTablesRef(linkRouteTableRef)

				mockOscRouteTableInterface.
					EXPECT().
					GetRouteTableIdsFromNetIds(gomock.Any(), gomock.Eq(netId)).
					Return([]string{routeTableId}, rttc.expGetRouteTableIdsFromNetIdsErr)

				routesSpec := routeTableSpec.Routes
				for _, routeSpec := range routesSpec {
					destinationIpRange := routeSpec.Destination
					resourceType := routeSpec.TargetType
					routeName := routeSpec.Name + "-uid"
					routeRef.ResourceMap[routeName] = routeTableId
					if resourceType == "gateway" {
						resourceId = internetServiceId
					} else {
						resourceId = natServiceId
					}
					route := osc.CreateRouteResponse{
						RouteTable: &osc.RouteTable{
							RouteTableId: &routeTableId,
						},
					}

					readRouteTables := osc.ReadRouteTablesResponse{
						RouteTables: &[]osc.RouteTable{
							*route.RouteTable,
						},
					}

					readRouteTable := *readRouteTables.RouteTables
					mockOscRouteTableInterface.
						EXPECT().
						GetRouteTableFromRoute(gomock.Any(), gomock.Eq(associateRouteTableId), gomock.Eq(resourceId), gomock.Eq(resourceType)).
						Return(&readRouteTable[0], rttc.expGetRouteTableFromRouteErr)
					mockOscRouteTableInterface.
						EXPECT().
						DeleteRoute(gomock.Any(), destinationIpRange, routeTableId).
						Return(rttc.expDeleteRouteErr)
				}
			}
			reconcileDeleteRouteTable, err := reconcileDeleteRouteTable(ctx, clusterScope, mockOscRouteTableInterface)
			if rttc.expReconcileDeleteRouteTableErr != nil {
				require.EqualError(t, err, rttc.expReconcileDeleteRouteTableErr.Error(), "reconcileDeleteRouteTable() should return the same error")
			} else {
				require.NoError(t, err)
			}
			t.Logf("Find reconcileDeleteRouteTable %v\n", reconcileDeleteRouteTable)
		})
	}
}

// TestReconcileDeleteRouteTableResourceId has several tests to cover the code of the function reconcileDeleteRouteTable
func TestReconcileDeleteRouteTableResourceId(t *testing.T) {
	routeTableTestCases := []struct {
		name                            string
		spec                            infrastructurev1beta1.OscClusterSpec
		expNetFound                     bool
		expReconcileDeleteRouteTableErr error
	}{
		{
			name: "check work without net, routeTable and route spec (with default values)",
			spec: infrastructurev1beta1.OscClusterSpec{
				Network: infrastructurev1beta1.OscNetwork{},
			},
			expNetFound:                     false,
			expReconcileDeleteRouteTableErr: errors.New("cluster-api-net-uid does not exist"),
		},
		{
			name:                            "net does not exist",
			spec:                            defaultRouteTableNatReconcile,
			expNetFound:                     false,
			expReconcileDeleteRouteTableErr: errors.New("test-net-uid does not exist"),
		},
	}
	for _, rttc := range routeTableTestCases {
		t.Run(rttc.name, func(t *testing.T) {
			clusterScope, ctx, mockOscRouteTableInterface, _ := SetupWithRouteTableMock(t, rttc.name, rttc.spec)

			netName := rttc.spec.Network.Net.Name + "-uid"
			netId := "vpc-" + netName

			netRef := clusterScope.GetNetRef()
			netRef.ResourceMap = make(map[string]string)
			if rttc.expNetFound {
				netRef.ResourceMap[netName] = netId
			}

			reconcileDeleteRouteTable, err := reconcileDeleteRouteTable(ctx, clusterScope, mockOscRouteTableInterface)
			if rttc.expReconcileDeleteRouteTableErr != nil {
				require.EqualError(t, err, rttc.expReconcileDeleteRouteTableErr.Error(), "reconcileDeleteRouteTable() should return the same error")
			} else {
				require.NoError(t, err)
			}
			t.Logf("Find reconcileDeleteRouteTable %v\n", reconcileDeleteRouteTable)
		})
	}
}
