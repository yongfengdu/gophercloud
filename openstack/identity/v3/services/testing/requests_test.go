package testing

import (
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateSuccessful(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateServiceSuccessfully(t)

	createOpts := services.CreateOpts{
		Type: "compute",
		Extra: map[string]interface{}{
			"name":        "service-two",
			"description": "Service Two",
			"email":       "service@example.com",
		},
	}

	actual, err := services.Create(client.ServiceClient(), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondService, *actual)
}

func TestListSinglePage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListServicesSuccessfully(t)

	count := 0
	err := services.List(client.ServiceClient(), services.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := services.ExtractServices(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedServicesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestGetSuccessful(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetServiceSuccessfully(t)

	actual, err := services.Get(client.ServiceClient(), "9876").Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondService, *actual)
	th.AssertEquals(t, SecondService.Extra["email"], "service@example.com")
}

func TestUpdateSuccessful(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateServiceSuccessfully(t)

	actual, err := services.Update(client.ServiceClient(), "9876", "compute2").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondServiceUpdated, *actual)
}

func TestDeleteSuccessful(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/services/12345", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := services.Delete(client.ServiceClient(), "12345")
	th.AssertNoErr(t, res.Err)
}
