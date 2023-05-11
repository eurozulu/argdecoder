package argdecoder

import (
	"fmt"
	"strings"
	"testing"
)

var testArgsNoFlags = strings.Split("one two three", " ")

var testArgsFollowingStrings = strings.Split("one two three -first 1st -second 2nd -third 3rd", " ")
var testArgsPreceedingStrings = strings.Split("-first 1st -second 2nd -third 3rd one two three", " ")
var testArgsMixedStrings = strings.Split("one -first 1st two -second 2nd three -third 3rd", " ")

var testArgsFollowingBool = strings.Split("one two three -first -second -third", " ")
var testArgsPreceedingBool = strings.Split("-first -second -third one two three", " ")
var testArgsMixedBool = strings.Split("one -first two -second three -third", " ")

type stringFlagStruct struct {
	First  string
	Second string
	Third  string
}

type stringPointerFlagStruct struct {
	First  *string
	Second *string
	Third  *string
}

type boolFlagStruct struct {
	First  bool
	Second bool
	Third  bool
}

func TestApplyArgumentsNoFlags(t *testing.T) {
	st := &stringFlagStruct{}
	args, err := ApplyArguments(testArgsNoFlags, st)
	if err != nil {
		t.Errorf("Failed to apply no flags argument  %v", err)
	}
	if !isSliceEqual(args, testArgsNoFlags) {
		t.Errorf("unexpected arguments returned applying no flags arguments.  Expected %v, found %v", testArgsNoFlags, args)
	}
	if st.First != "" || st.Second != "" || st.Third != "" {
		t.Errorf("unexpected field value assigned with no flag arguments")
	}
}

func TestApplyArgumentsStringFlags(t *testing.T) {
	err := testStringFlags(testArgsFollowingStrings)
	if err != nil {
		t.Errorf("Failed to apply string flags for testArgsFollowingStrings argument  %v", err)
	}
	err = testStringFlags(testArgsPreceedingStrings)
	if err != nil {
		t.Errorf("Failed to apply string flags for testArgsPreceedingStrings argument  %v", err)
	}
	err = testStringFlags(testArgsMixedStrings)
	if err != nil {
		t.Errorf("Failed to apply string flags for testArgsMixedStrings argument  %v", err)
	}
}

func TestApplyArgumentsStringPointerFlags(t *testing.T) {
	err := testStringPointerFlags(testArgsFollowingStrings)
	if err != nil {
		t.Errorf("Failed to apply string pointer flags for testArgsFollowingStrings argument  %v", err)
	}
	err = testStringPointerFlags(testArgsPreceedingStrings)
	if err != nil {
		t.Errorf("Failed to apply string pointer flags for testArgsPreceedingStrings argument  %v", err)
	}
	err = testStringPointerFlags(testArgsMixedStrings)
	if err != nil {
		t.Errorf("Failed to apply string pointer flags for testArgsMixedStrings argument  %v", err)
	}

	// test nil value
	st := &stringPointerFlagStruct{}
	if _, err = ApplyArguments([]string{"-first", "-second", "2nd"}, st); err != nil {
		t.Errorf("failed to apply nil value test pointers  %v", err)
		return
	}
	if st.First != nil {
		t.Errorf("unexpected value found in First where nil expected.  Found '%v'", *st.First)
	}
	if st.Third != nil {
		t.Errorf("unexpected value found in Third where nil expected.  Found %v", st.Third)
	}
	if st.Second == nil {
		t.Errorf("unexpected nil value found in Second. expected %s", "2nd")
	} else if *st.Second != "2nd" {
		t.Errorf("unexpected value found in Second. expected %s, found %s", "2nd", *st.Second)
	}
}

func TestApplyArgumentsBoolFlags(t *testing.T) {
	err := testBoolFlags(testArgsFollowingBool)
	if err != nil {
		t.Errorf("Failed to apply string flags for testArgsFollowingStrings argument  %v", err)
	}
	err = testBoolFlags(testArgsPreceedingBool)
	if err != nil {
		t.Errorf("Failed to apply string flags for testArgsPreceedingStrings argument  %v", err)
	}
	err = testBoolFlags(testArgsMixedBool)
	if err != nil {
		t.Errorf("Failed to apply string flags for testArgsMixedStrings argument  %v", err)
	}
}

func testStringFlags(args []string) error {
	st := &stringFlagStruct{}
	remain, err := ApplyArguments(args, st)
	if err != nil {
		return fmt.Errorf("Failed to apply no argument  %v", err)
	}
	if !isSliceEqual(remain, testArgsNoFlags) {
		return fmt.Errorf("unexpected remaining arguments applying string flags.  Expected %v, found %v", testArgsNoFlags, args)
	}
	if st.First != "1st" {
		return fmt.Errorf("unexpected field value in First. Expected %s, found %s", "1st", st.First)
	}
	if st.Second != "2nd" {
		return fmt.Errorf("unexpected field value in Second. Expected %s, found %s", "2nd", st.Second)
	}
	if st.Third != "3rd" {
		return fmt.Errorf("unexpected field value in Third. Expected %s, found %s", "3rd", st.Third)
	}
	return nil
}

func testStringPointerFlags(args []string) error {
	st := &stringPointerFlagStruct{}
	remain, err := ApplyArguments(args, st)
	if err != nil {
		return fmt.Errorf("Failed to apply no argument  %v", err)
	}
	if !isSliceEqual(remain, testArgsNoFlags) {
		return fmt.Errorf("unexpected remaining arguments applying string flags.  Expected %v, found %v", testArgsNoFlags, args)
	}
	if st.First == nil {
		return fmt.Errorf("unexpected field value in First. Expected %s, found nil", "1st")
	}
	if *st.First != "1st" {
		return fmt.Errorf("unexpected field value in First. Expected %s, found %s", "1st", *st.First)
	}
	if st.Second == nil {
		return fmt.Errorf("unexpected field value in Second. Expected %s, found nil", "2nd")
	}
	if *st.Second != "2nd" {
		return fmt.Errorf("unexpected field value in Second. Expected %s, found %s", "2nd", *st.Second)
	}
	if st.Third == nil {
		return fmt.Errorf("unexpected field value in Third. Expected %s, found nil", "3rd")
	}
	if *st.Third != "3rd" {
		return fmt.Errorf("unexpected field value in Third. Expected %s, found %s", "3rd", *st.Third)
	}

	return nil
}

func testBoolFlags(args []string) error {
	st := &boolFlagStruct{}
	remain, err := ApplyArguments(args, st)
	if err != nil {
		return fmt.Errorf("Failed to apply no argument  %v", err)
	}
	if !isSliceEqual(remain, testArgsNoFlags) {
		return fmt.Errorf("unexpected remaining arguments applying bool flags.  Expected %v, found %v", testArgsNoFlags, args)
	}
	if !st.First {
		return fmt.Errorf("unexpected field value in First. Expected %v, found %v", true, st.First)
	}
	if !st.Second {
		return fmt.Errorf("unexpected field value in Second. Expected %v, found %v", true, st.Second)
	}
	if !st.Third {
		return fmt.Errorf("unexpected field value in Third. Expected %v, found %v", true, st.Third)
	}
	return nil
}

func isSliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, ss1 := range s1 {
		if ss1 != s2[i] {
			return false
		}
	}
	return true
}
