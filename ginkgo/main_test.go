package main_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bclicn/color"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/salihkemaloglu/UnitAndIntegrationTesting-Golang/operations"
)

/*TODO:
Error handling proglem
Docker deploy problem
*/
var baseUrl = "http://api:8080/item"

var _ = Describe(color.LightBlue("Integration test Webapi GetAll"), func() {
	Context(color.LightGreen("initially"), func() {
		response := GetAll()
		itemGet := response[len(response)-1]
		It(color.LightYellow("has 0 items"), func() {
			Expect("Nemo").Should(Equal(itemGet.Name))
		})
		Specify(color.LightYellow("the total item is 0"), func() {
			Expect(1).Should(Equal(len(response)))
		})
	})
})
var _ = Describe(color.LightBlue("Integration test Webapi GetAll two Decribe in Describe with different Contexts"), func() {
	response := GetAll()
	itemGet := response[len(response)-1]
	Describe(color.LightBlue("Testing Describe inside describe"), func() {
		Context(color.LightGreen("Context 1"), func() {
			It(color.LightYellow("has 0 items"), func() {
				Expect("Nemo").Should(Equal(itemGet.Name))
			})
		})
		Context(color.LightGreen("Context 2"), func() {
			Specify(color.LightYellow("the total item is 0"), func() {
				Expect(1).Should(Equal(len(response)))
			})
		})
	})
})

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration test Webapi")
}
func GetAll() []Item {
	response, err := http.Get(baseUrl)
	if err != nil {
		errMessage := color.LightYellow("End point does not responde!:" + err.Error())
		Specify(errMessage, func() {
			Expect(nil).Should(Equal(err.Error()))
		})
		return nil
	} else if response.StatusCode != 200 {
		errMessage := color.LightYellow("Server side response not Ok!,Response StatusCode:" + string(response.StatusCode))
		Specify(errMessage, func() {
			Expect(nil).Should(Equal(err))
		})
		return nil
	} else {
		defer response.Body.Close()
		var item []Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			errMessage := color.LightYellow("Json decode error!:" + err.Error())
			Specify(errMessage, func() {
				Expect(nil).Should(Equal(err))
			})
			return nil
		} else {
			return item
		}
	}
}
