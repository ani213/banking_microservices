package domain

type CustomerRepositoryStub struct {
	customer []Customer
}

func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customer, nil
}

func NewCustomerRepositoryStub() CustomerRepositoryStub {
	customer := []Customer{
		{Id: "1", Name: "John", City: "New York", ZipCode: "10001", DateOfBirth: "1990-01-01", Status: "active"},
		{Id: "2", Name: "Jane", City: "Los Angeles", ZipCode: "90001", DateOfBirth: "1992-02-02", Status: "inactive"},
		{Id: "3", Name: "Bob", City: "Chicago", ZipCode: "60601", DateOfBirth: "1985-03-03", Status: "active"},
	}
	return CustomerRepositoryStub{customer: customer}
}
