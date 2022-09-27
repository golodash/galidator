CHANGELOG
=========

UNRELEASED
----------

* 🎉 feat: validator creation from a slice is possible
* ✅ test: FromStruct testcase added
* ✅ test: a complex testcase added
* ✅ test: WhenExistAll testcases added
* ✅ test: WhenExistOne testcases added
* ✅ test: Type testcases added
* ✅ test: String testcases added
* ✅ test: Choices testcases added
* ✅ test: XOR testcases added
* ✅ test: OR testcases added
* ✅ test: Password testcases added
* ✅ test: Slice testcases added
* ✅ test: Struct testcases added
* ✅ test: Map testcases added
* ✅ test: Phone testcases added
* ✅ test: Regex testcases added
* ✅ test: Email testcases added
* ✅ test: NonEmpty testcases added
* ✅ test: NonNil testcases added
* ✅ test: NonZero testcases added
* ✅ test: Required testcases added
* ✅ test: Len testcases added
* ✅ test: LenRange testcases added
* ✅ test: max testcases added
* 🐛 fix: specificMessages and messages are fixed
* ✅ test: min testcases added
* ✅ test: float testcases added
* ✅ test: int test cases added
* 🎉 feat: OR and XOR is possible inside tags now
* 🐛 fix: whenExistAllRequireRule bug fixed
* 🎉 feat: Choices, WhenExistOne and WhenExistAll are possible by tags now
* 🎉 feat: adding specificMessages inside a ruleSet is possible now
* 🎉 feat: custom functions can be defined by elements tags
* 🐛 fix: Messages and SpecificMessages work just fine again
* 🎉 feat: added Validator creation based on structs elements tags
* 🐛 fix: validation bug fixed on slice types being required
* 🎉 feat: fully flexible validator structure
* 🎉 feat: added slice validation + more complex validations are possible
* 🎉 feat: String rule added + default errors can be set when generating generator + struct validation added
* 🎉 feat: added a name for customizing output of fields
* 🐛 fix: WhenExistAll bug fixed
* 🎉 feat: added WhenExistAll rule function, added Optional rule function
* 🎉 feat: added WhenExistAll rule function, added Optional function
* 🎉 feat: requires added core structure of library
* 🎉 feat: choices function added
* 🐛 fix: return nill when no errors happened
* 🎉 feat: XOR rule added, Or changed to OR
* 🎉 feat: Or rule added
* 🎉 feat: added R for RuleSet
* 🐛 fix: Pointer name does not exists in 1.17, so back to reflect.Ptr name which is supported by v1.17 of golang
* 🎉 feat: if pass as reference happens, Validate function now can handle it too
* 🐛 fix: fixed error on adding custom validators
* 🐛 fix: bug fixed on elements of slice
* 🎉 feat: Password validator added
* 🐛 fix: slice validation children added
* 🎉 feat: slice elements with complex can be checked
* 🐛 fix: Validate function returns nil if no error happened
* 🎉 feat: $fieldS addes snake_case version of $field to the error output message
* 🎉 feat: struct and map can be used
* 🐛 fix: required will happen in other function in some cases
* 🎉 feat: custom function added
* 🎉 feat: phone validator added
* 🎉 feat: ability to add custom error messages for specific rules added in specific fields added
* 🎉 feat: Regex function added
feat: error output on error validations fixed
* 🎉 feat: all fields are optional, until developer uses Required, NonNil, NonEmpty or NonZero rules
feat: Email validator added
* 🎉 feat: NonEmpty function added
* 🎉 feat: added NonNil and NonZero
* 🐛 fix: making Min, Max, LenRange and Len general usecase
* 🎉 feat: added Len and LenRange functions in different meanings
* 🐛 fix: workflow now downloads requirements
* 🐛 fix: Required paniced when int values been used
* 🎉 feat: Required function added
* 🎉 feat: users can now pass rules' keys as PascalCase or snake_case
* 🎉 feat: Len validator added
* 🎉 feat(validators): Min and Max function and options parameter added to use in error prints
* 🎉 feat(utils): determinePrecision function added

.. 1.0.0 (2022-06-22)
.. ------------------
