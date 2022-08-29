package VkApiTest

import (
	"errors"
	"testing"
	VkApi "vk-api-tools"
)

func TestUseCase(t *testing.T) {
	any := VkApi.AnyModel(`{
  "first_name":"Ivan",
  "last_name":"Ivanov",
  "city":{
    "id":19,
    "title":"Moscow"
  },
  "friends":[ {"id":1,"name":"Jon","verify":true},{"id":2,"name":"Bob","verify":false},{"id":3,"name":"Alexandra"} ]
}`)

	firstName := any.QStringDef("first_name", "")
	cityName := any.QStringDef("city.title", "")

	hasAge := any.QString("age") != nil

	firstFriendName := any.QStringDef("friends.0.name", "")
	lastFriendId := any.QIntDef("friends.-1.id", -1)

	if firstName != "Ivan" {
		t.Error("Expexted Ivan, got", firstName)
	}

	if cityName != "Moscow" {
		t.Error("Expected Moscow, got", cityName)
	}

	if firstFriendName != "Jon" {
		t.Error("Expected Jon, got", firstFriendName)
	}

	if lastFriendId != 3 {
		t.Error("Expected 3, got", lastFriendId)
	}

	if hasAge != false {
		t.Error("Expected no age, got", hasAge)
	}

	friends := any.QSlice("friends")

	buff := ""
	if friends != nil {
		for _, friend := range *friends {
			v := friend.QBool("verify")

			if v == nil || *v == true {
				buff = buff + friend.QStringDef("name", "") + "\n"
			}
		}
	}

	if buff != "Jon\nAlexandra\n" {
		t.Error("Expected, Jon\nAlexandra\n, got", buff)
	}

}

func TestAnyModelBasicInt(t *testing.T) {
	any := VkApi.AnyModel("10")

	i, err := any.Int()

	if err != nil {
		t.Error(err)
	}

	if i != 10 {
		t.Error("Expected 10, got", i)
	}
}

func TestAnyModelBasicBool(t *testing.T) {
	any := VkApi.AnyModel("true")
	i, err := any.Bool()

	if err != nil {
		t.Error(err)
	}

	if i != true {
		t.Error("Expected true, got", i)
	}

	any = VkApi.AnyModel("false")

	i, err = any.Bool()

	if err != nil {
		t.Error(err)
	}

	if i != false {
		t.Error("Expected false, got", i)
	}

	any = VkApi.AnyModel("krokodail")

	i, err = any.Bool()

	if err == nil {
		t.Error("Expected error, but not error")
	}
}

func TestAnyModelBasicFloat(t *testing.T) {
	any := VkApi.AnyModel("5")
	i, err := any.Float()

	if err != nil {
		t.Error(err)
	}

	if i != 5 {
		t.Error("Expected 5, got", i)
	}

	any = VkApi.AnyModel("5.666")

	i, err = any.Float()

	if err != nil {
		t.Error(err)
	}

	if i != 5.666 {
		t.Error("Expected 5.666, got", i)
	}

	any = VkApi.AnyModel("krokodail")

	i, err = any.Float()

	if err == nil {
		t.Error("Expected error, but not error")
	}
}

func TestAnyModelBasicAny(t *testing.T) {
	any := VkApi.AnyModel(`{"name":"Ivan","age":15}`)
	name, err := any.GetAny("name")

	if err != nil {
		t.Error(err)
	}

	if s, _ := name.String(); s != "Ivan" {
		t.Error("Expected Ivan, got", s)
	}
}

func TestAnyModelBasicString(t *testing.T) {
	any := VkApi.AnyModel("\"10\"")

	i, err := any.String()

	if err != nil {
		t.Error(err)
	}

	if i != "10" {
		t.Error("Expected \"10\", got", i)
	}
}

func TestAnyModelBasicGetStringByKey(t *testing.T) {
	any := VkApi.AnyModel(`{"one":"1","two":"2","name":"Ivan"}`)

	if one, err := any.GetString("one"); one != "1" {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Expected 1, got", one)
		}
	}

	if one, err := any.GetStringDef("one", "1"); one != "1" {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Expected 1, got", one)
		}
	}

	if one, err := any.GetStringDef("not_exist_key", "1"); one != "1" {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Expected 1, got", one)
		}
	}

	if one, err := any.GetString("name"); one != "Ivan" {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Expected Ivan, got", one)
		}
	}
}

func TestAnyModelBasicGetIntByKey(t *testing.T) {
	any := VkApi.AnyModel(`{"one":1,"two":2,"name":"Ivan"}`)

	if one, err := any.GetInt("one"); one != 1 {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Expected 1, got", one)
		}
	}

	if one, err := any.GetIntDef("one", 1); one != 1 {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Expected 1, got", one)
		}
	}

	if one, err := any.GetIntDef("not_exist_key", 1); one != 1 {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Expected 1, got", one)
		}
	}

	if one, err := any.GetInt("two"); one != 2 {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Expected 2, got", one)
		}
	}
}

func tss(query string, result string, any VkApi.AnyModel) error {
	last := any.QString(query)
	if last == nil {
		return errors.New("Expected " + result + ", got nil on query " + query)
	} else if *last != result {
		errors.New("Expected " + result + ", got" + *last + " on query " + query)
	}
	return nil
}

func TestAnyModelQuery(t *testing.T) {
	any := VkApi.AnyModel(`{
"int":199,
"interstring":"199",
"string":"Hello world",
"bool":false,
"array":[199,"Hello world",false,true,["Hello world", 199],{"int":99,"string":"Hello world"}],
"object":{"0":199,"1":"Hello world","array":[199,"Hello world",199],"bool":false}
}`)

	stringQueries := []string{
		"string",
		"array.1",
		"array.0",
		"array.-5",
		"array.4.0",
		"array.4.1",
		"array.4.-2",
		"array.5.string",
		"object.1",
		"object.array.1",
		"object.array.-2",
	}

	intQueries := []string{
		"int",
		"array.0",
		"array.4.1",
		"object.0",
		"interstring",
	}

	boolQueries := []string{
		"bool",
		"array.2",
		"object.bool",
	}

	var e error
	for _, q := range stringQueries {
		e = tss(q, "Hello world", any)
		if e != nil {
			t.Error(e)
		}
	}

	if any.QString("object.array") != nil {
		t.Error("Expected object.array is null, got string")
	}

	if any.QStringDef("object.array", "Hello") != "Hello" {
		t.Error("Expected object.array is def Hello, got another string")
	}

	if any.QStringDef("array.5.string", "Hello") != "Hello world" {
		t.Error("Expected object.array is Hello world, got another string")
	}

	for _, q := range intQueries {
		res := any.QInt(q)
		if res == nil {
			t.Error("Expected 199 at query", q, "got nil")
		} else if *res != 199 {
			t.Error("Expected 199 at query", q, "got", *res)
		}
	}

	for _, q := range boolQueries {
		res := any.QBool(q)
		if res == nil {
			t.Error("Expected false at query", q, "got nil")
		} else if *res != false {
			t.Error("Expected false at query", q, "got", *res)
		}
	}
}
