package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGetStringBodyHttpRequestJSON(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"foo": "bar"})
	req, _ := http.NewRequest("POST", "http://server.com", bytes.NewBuffer(body))
	actual := GetStringBodyHttpRequest(req)

	assert.NotNil(t, actual)
	assert.Equal(t, "{\"foo\":\"bar\"}", *actual)
}

func TestGetStringBodyHttpRequestPlainText(t *testing.T) {
	stringBody := "PLAIN TEXT"
	byteArrayStringBody := []byte(stringBody)
	req, _ := http.NewRequest("POST", "http://server.com", bytes.NewBuffer(byteArrayStringBody))
	actual := GetStringBodyHttpRequest(req)

	assert.NotNil(t, actual)
	assert.Equal(t, stringBody, *actual)
}

func TestGetStringBodyHttpRequestJSONEncoded(t *testing.T) {
	stringBody := `1223ab
{'response':{'code':200}}
0

`
	byteArrayStringBody := []byte(stringBody)
	req, _ := http.NewRequest("POST", "http://server.com", bytes.NewBuffer(byteArrayStringBody))
	actual := GetStringBodyHttpRequestJSON(req)

	assert.NotNil(t, actual)
	assert.Equal(t, "{'response':{'code':200}}", *actual)
}

func TestGetStringBodyHttpResponseJSON(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	defer gock.Clean()

	gock.New("http://server.com").
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	req, _ := http.NewRequest("GET", "http://server.com/bar", nil)
	client := &http.Client{}
	res, _ := client.Do(req)
	actual := GetStringBodyHttpResponse(res)

	assert.NotNil(t, actual)
	assert.Equal(t, "{\"foo\":\"bar\"}\n", *actual)
}

func TestGetStringBodyHttpResponsePlainText(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	defer gock.Clean()

	stringBody := "PLAIN TEXT"

	gock.New("http://server.com").
		Get("/bar").
		Reply(200).
		BodyString(stringBody)

	req, err := http.NewRequest("GET", "http://server.com/bar", nil)
	if err != nil {
		fmt.Printf("\nErr: %+v\n\n", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nErr: %+v\n\n", err)
	}
	actual := GetStringBodyHttpResponse(res)

	assert.NotNil(t, actual)
	assert.Equal(t, stringBody, *actual)
}

func TestGetStringBodyHttpResponseJSONEncoded(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	defer gock.Clean()

	stringBody := `1223ab
{'response':{'code':200}}
0

`

	gock.New("http://server.com").
		Get("/bar").
		Reply(200).
		BodyString(stringBody)

	req, err := http.NewRequest("GET", "http://server.com/bar", nil)
	if err != nil {
		fmt.Printf("\nErr: %+v\n\n", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nErr: %+v\n\n", err)
	}
	actual := GetStringBodyHttpResponseJSON(res)

	assert.NotNil(t, actual)
	assert.Equal(t, "{'response':{'code':200}}", *actual)
}

func TestIntegerSliceToStringSlice(t *testing.T) {
	strs := ToStringSlice([]int{1, 2, 3})

	assert.Len(t, strs, 3)
	assert.Equal(t, "1", strs[0])
	assert.Equal(t, "2", strs[1])
	assert.Equal(t, "3", strs[2])
}

func TestInteger64SliceToStringSlice(t *testing.T) {
	strs := ToStringSlice64([]int64{16, 23, 39})

	assert.Len(t, strs, 3)
	assert.Equal(t, "16", strs[0])
	assert.Equal(t, "23", strs[1])
	assert.Equal(t, "39", strs[2])
}

func TestToInt64Slice(t *testing.T) {
	actual := ToInt64Slice([]string{"654987", "852369", "a"})

	assert.Len(t, actual, 2)
	assert.Equal(t, int64(654987), actual[0])
	assert.Equal(t, int64(852369), actual[1])
}

func TestParseInt(t *testing.T) {
	i, err := ParseStringToInt("6549")

	expected := int(6549)

	assert.Empty(t, err)
	assert.IsType(t, expected, i)
	assert.Equal(t, expected, i)
}

func TestParseIntWithEmptyString(t *testing.T) {
	i, err := ParseStringToInt("")

	assert.Equal(t, 0, i)
	assert.Empty(t, err)
}

func TestParseIntInvalidString(t *testing.T) {
	_, err := ParseStringToInt("invalid")

	assert.NotEmpty(t, err)
}

func TestParseInt64(t *testing.T) {
	i, err := ParseStringToInt64("456123789123")

	expected := int64(456123789123)

	assert.Empty(t, err)
	assert.IsType(t, expected, i)
	assert.Equal(t, expected, i)
}

func TestParseInt64WithEmptyString(t *testing.T) {
	i, err := ParseStringToInt64("")

	assert.Equal(t, int64(0), i)
	assert.Empty(t, err)
}

func TestParseInt64InvalidString(t *testing.T) {
	_, err := ParseStringToInt64("invalid")

	assert.NotEmpty(t, err)
}

func TestShouldParseTimeWithYearMonthDayPattern(t *testing.T) {
	date, err := ParseDateYearMonthDay("2000-12-31")
	assert.Nil(t, err)
	assert.False(t, date.IsZero())
	assert.EqualValues(t, 2000, date.Year())
	assert.EqualValues(t, 12, date.Month())
	assert.EqualValues(t, 31, date.Day())
	assert.EqualValues(t, 0, date.Hour())
	assert.EqualValues(t, 0, date.Minute())
	assert.EqualValues(t, 0, date.Second())
}

func TestShouldNotParseTimeWithoutYearMonthDayPattern(t *testing.T) {
	var err error
	_, err = ParseDateYearMonthDay("01-12-2000")
	assert.NotNil(t, err)

	_, err = ParseDateYearMonthDay("01-12-00")
	assert.NotNil(t, err)
}

func TestDiffDays(t *testing.T) {
	duration, err := DiffDays(time.Date(2016, 2, 5, 0, 0, 0, 0, time.UTC), time.Date(2016, 2, 11, 0, 0, 0, 0, time.UTC))
	assert.NotNil(t, duration)
	assert.Nil(t, err)
	assert.Equal(t, int64(6), duration)

	duration, err = DiffDays(time.Date(2016, 2, 20, 0, 0, 0, 0, time.UTC), time.Date(2016, 3, 10, 0, 0, 0, 0, time.UTC))
	assert.NotNil(t, duration)
	assert.Nil(t, err)
	assert.Equal(t, int64(19), duration)

	date1 := time.Time{}
	date2 := time.Time{}
	duration, err = DiffDays(date1, date2)
	assert.Empty(t, duration)
	assert.NotNil(t, err)
}

func TestShouldParseDateStringMalformedTimeToTime(t *testing.T) {
	var expected time.Time

	expected, _ = time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")

	result1, err := ParseDateStringToTime("2016-01-01")
	assert.Nil(t, err)
	assert.Equal(t, expected, *result1)

	result2, err := ParseDateStringToTime("2016-01-01T00:00:00")
	assert.Nil(t, err)
	assert.Equal(t, expected, *result2)

	result3, err := ParseDateStringToTime("2016-01-01T00:00:00Z")
	assert.Nil(t, err)
	assert.Equal(t, expected, *result3)

	result4, err := ParseDateStringToTime("2016-01-01 00:00:00")
	assert.Nil(t, err)
	assert.Equal(t, expected, *result4)

	_, err = ParseDateStringToTime("2016-01-01T00:00:00ABC")
	assert.NotNil(t, err)
}

func TestShouldParseDateStringMalformedTimeToTimeZero(t *testing.T) {
	expected := time.Time{}

	result1, err := ParseDateStringToTime("0000-00-00")
	assert.Nil(t, err)
	assert.Equal(t, expected, *result1)

	result2, err := ParseDateStringToTime("0000-00-00T00:00:00")
	assert.Nil(t, err)
	assert.Equal(t, expected, *result2)

	result3, err := ParseDateStringToTime("0000-00-00T00:00:00Z")
	assert.Nil(t, err)
	assert.Equal(t, expected, *result3)

	result4, err := ParseDateStringToTime("0000-00-00 00:00:00")
	assert.Nil(t, err)
	assert.Equal(t, expected, *result4)

	_, err = ParseDateStringToTime("0000-00-00T00:00:00ABC")
	assert.NotNil(t, err)
}

func TestShouldParseIntToBool(t *testing.T) {
	result1 := ParseIntToBool(0)
	assert.Equal(t, false, result1)

	result2 := ParseIntToBool(1)
	assert.Equal(t, true, result2)

	result3 := ParseIntToBool(2)
	assert.Equal(t, false, result3)

	result4 := ParseIntToBool(345)
	assert.Equal(t, false, result4)
}

func TestShouldParseBoolToString(t *testing.T) {
	result1 := ParseBoolToString(true)
	assert.Equal(t, "1", result1)

	result2 := ParseBoolToString(false)
	assert.Equal(t, "0", result2)
}

func TestShouldCheckStringJsonData(t *testing.T) {
	var s string
	result1 := CheckStringJsonData(s)
	assert.Nil(t, result1)

	result2 := CheckStringJsonData("")
	assert.Nil(t, result2)

	result3 := CheckStringJsonData("test")
	assert.NotNil(t, result3)
	assert.Equal(t, "test", *result3)
}

func TestShouldCheckInt64JsonData(t *testing.T) {
	var i1 int64
	result1 := CheckInt64JsonData(i1)
	assert.Nil(t, result1)

	result2 := CheckInt64JsonData(0)
	assert.Nil(t, result2)

	result3 := CheckInt64JsonData(987654)
	assert.NotNil(t, result3)
	assert.Equal(t, int64(987654), *result3)
}

func TestShouldCheckFloat64JsonData(t *testing.T) {
	var f1 float64
	result1 := CheckFloat64JsonData(f1)
	assert.Nil(t, result1)

	result2 := CheckFloat64JsonData(0)
	assert.Nil(t, result2)

	result3 := CheckFloat64JsonData(0.00)
	assert.Nil(t, result3)

	result4 := CheckFloat64JsonData(9876.54)
	assert.NotNil(t, result4)
	assert.Equal(t, float64(9876.54), *result4)
}

func TestShouldReturnOnlyNumbers(t *testing.T) {
	s1 := "61.225.412/0001-14aA"

	result := *GetOnlyNumbers(&s1)
	assert.Equal(t, "61225412000114", result)
}

func TestShouldReturnOnlyNumbersOrSpecial(t *testing.T) {
	s1 := "+55 (21) 98765-4321"

	result := *GetOnlyNumbersOrSpecial(&s1, "+")
	assert.Equal(t, "+5521987654321", result)
}

func TestShouldReturnOnlyNumbersOrSpecial1(t *testing.T) {
	s1 := "+55 (21) 98765-4321"

	result := *GetOnlyNumbersOrSpecial(&s1, "+()")
	assert.Equal(t, "+55(21)987654321", result)
}

func TestShouldReturnNilForNilInput(t *testing.T) {
	var s1 string

	result := *GetOnlyNumbers(&s1)
	assert.Equal(t, s1, result)
}

func TestShouldReturnNilForNilInput1(t *testing.T) {
	var s1 string

	result := *GetOnlyNumbersOrSpecial(&s1, "+")
	assert.Equal(t, s1, result)
}

func TestParseIntOrReturnZero(t *testing.T) {
	stg := "1"
	expected := 1

	assert.Equal(t, expected, ParseIntOrReturnZero(stg))
}

func TestParseIntOrReturnZeroFail(t *testing.T) {
	stg := "a"
	expected := 0

	assert.Equal(t, expected, ParseIntOrReturnZero(stg))
}

func TestParseIntOrReturnZeroWithNumberOnString(t *testing.T) {
	stg := "a123"
	expected := 0

	assert.Equal(t, expected, ParseIntOrReturnZero(stg))
}

func TestIsArray(t *testing.T) {
	actual := IsArray([]string{"foo", "bar"})
	assert.Equal(t, true, actual)

	actual = IsArray([]int{65485, 19734})
	assert.Equal(t, true, actual)

	actual = IsArray([]int64{65485, 19734})
	assert.Equal(t, true, actual)

	actual = IsArray(nil)
	assert.Equal(t, false, actual)

	actual = IsArray(65485)
	assert.Equal(t, false, actual)

	actual = IsArray("foo")
	assert.Equal(t, false, actual)

	actual = IsArray(false)
	assert.Equal(t, false, actual)
}

func TestJoin(t *testing.T) {
	actual := Join(", ", 654321987, "bar", 654.654)
	assert.Equal(t, `654321987, bar, 654.654`, actual)

	actual = Join(", ", int64(654321987), "bar")
	assert.Equal(t, `654321987, bar`, actual)

	actual = Join(", ", int64(654321987), int64(52354))
	assert.Equal(t, `654321987, 52354`, actual)

	actual = Join(", ", "foo")
	assert.Equal(t, `foo`, actual)

	actual = Join(", ", []string{"foo", "bar"})
	assert.Equal(t, `foo, bar`, actual)

	actual = Join(", ", []int{65485, 19734})
	assert.Equal(t, `65485, 19734`, actual)
}

func TestBeginningOfToday(t *testing.T) {
	today := BeginningOfToday()
	assert.Equal(t, today.Year(), time.Now().Year())
	assert.Equal(t, today.Month(), time.Now().Month())
	assert.Equal(t, today.Day(), time.Now().Day())
	assert.Equal(t, today.Hour(), 0)
	assert.Equal(t, today.Minute(), 0)
	assert.Equal(t, today.Second(), 0)
}

func TestShouldRemoveNanoseconds(t *testing.T) {
	expected := time.Date(2016, time.September, 20, 18, 49, 15, 0, time.UTC)

	date := time.Date(2016, time.September, 20, 18, 49, 15, 999999999, time.UTC)

	actual, err := RemoveNanoseconds(date)
	assert.Nil(t, err)
	assert.EqualValues(t, expected.Year(), actual.Year())
	assert.EqualValues(t, expected.Month(), actual.Month())
	assert.EqualValues(t, expected.Day(), actual.Day())
	assert.EqualValues(t, expected.Hour(), actual.Hour())
	assert.EqualValues(t, expected.Minute(), actual.Minute())
	assert.EqualValues(t, expected.Second(), actual.Second())
	assert.EqualValues(t, expected.Local(), actual.Local())
	assert.EqualValues(t, expected.Nanosecond(), actual.Nanosecond(), "need to be zeroed")
}