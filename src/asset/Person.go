package asset

import (
	"fmt"
	"time"
)

type Person struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,gte=0,lte=130"`
}

type Doctor struct {
	Person     `json:"person"`
	Department string `json:"department" validate:"required"`
}

type OutPatient struct {
	Person    `json:"person"`
	Country   string `json:"country"`
	Region    string `json:"region"`
	Birthday  string `json:"birthday" validate:"required,datetime=2006-01-02"`
	IsMarried bool   `json:"isMarried"`
	Career    string `json:"career"`
	Address   string `json:"address"`
}

/**
 * Getters And Setters
 */
func (p Person) GetName() string {
	return p.Name
}

func (p Person) SetName(name string) {
	p.Name = name
}

func (p Person) GetID() string {
	return p.ID
}

func (p Person) SetID(id string) {
	p.ID = id
}

func (p Person) GetAge() int {
	return p.Age
}

func (p Person) SetAge(age int) {
	p.Age = age
}

func (d Doctor) GetDepartment() string {
	return d.Department
}

func (d Doctor) SetDepartment(dep string) {
	d.Department = dep
}

func (o OutPatient) GetBirthday() time.Time {
	todayZero, _ := time.ParseInLocation("2006-01-02 15:04:05", o.Birthday, time.Local)
	return todayZero
}

func (o OutPatient) SetBirthday(birth time.Time) {
	o.Birthday = birth.Format("2006-01-02 15:04:05")
}

func (o OutPatient) GetCountry() string {
	return o.Country
}

func (o OutPatient) SetCountry(country string) {
	o.Country = country
}

func (o OutPatient) GetRegion() string {
	return o.Region
}

func (o OutPatient) SetRegion(region string) {
	o.Region = region
}

func (o OutPatient) Married() bool {
	return o.IsMarried
}

func (o OutPatient) SetMarried(m bool) {
	o.IsMarried = m
}

func (o OutPatient) GetCareer() string {
	return o.Career
}

func (o OutPatient) SetCareer(career string) {
	o.Career = career
}

func (o OutPatient) GetAddress() string {
	return o.Address
}

func (o OutPatient) SetAddress(addr string) {
	o.Address = addr
}

/**
 * ToString()s
 */

func (p Person) String() string {
	return fmt.Sprintf("<Person> id=%s,name=%s", p.ID, p.Name)
}

func (d Doctor) String() string {
	return fmt.Sprintf("<<Doctor %s> dep=%s>", d.Person.String(), d.Department)
}

func (o OutPatient) String() string {
	return fmt.Sprintf("<<OutPatient %s> isMarried=%t>", o.Person.String(), o.IsMarried)
}
