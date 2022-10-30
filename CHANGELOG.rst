CHANGELOG
=========

UNRELEASED
----------

* 🐛 fix: two bugs fixed
1. Capital letters in tags didn't register right in validator system for error SpecialMessages
2. More than one special error message in tags, just kept the last one in memory

0.0.0 (2022-09-29)
------------------

* ✅ test: test cases for many scenarios added
* 🐛 fix: Huge amount of bug fixes happened
* 🎉 feat: validator creation from a slice is possible
* 🎉 feat: OR and XOR is possible inside tags
* 🎉 feat: Choices, WhenExistOne and WhenExistAll are possible by tags
* 🎉 feat: adding specificMessages inside a ruleSet is possible
* 🎉 feat: custom functions can be defined by elements tags
* 🎉 feat: added slice validation + more complex validations like complex inside them is possible
* 🎉 feat: all fields are optional, until developer uses Required, NonNil, NonEmpty or NonZero rules or WhenExistAll or WhenExistOne
* 🎉 feat: String rule added + struct validation added
* 🎉 feat: added a name for ruleSets to customize output of fields
* 🎉 feat: added WhenExistAll rule function + added Optional rule function
* 🎉 feat: requires functionality added to core structure of library
* 🎉 feat: choices function added
* 🎉 feat: XOR rule added
* 🎉 feat: OR rule added
* 🎉 feat: added R as an alias for RuleSet
* 🎉 feat: if pass as reference happens, Validate function now can handle it
* 🎉 feat: Password validator added
* 🎉 feat: struct and map can be used as input
* 🎉 feat: custom function added
* 🎉 feat: phone validator added
* 🎉 feat: ability to add custom error messages for specific rules in specific fields added
* 🎉 feat: Regex function added
* 🎉 feat: Email validator added
* 🎉 feat: NonEmpty function added
* 🎉 feat: added NonNil and NonZero
* 🎉 feat: added Len and LenRange functions in different meanings
* 🎉 feat: Required function added
* 🎉 feat: users can now pass rules' keys as PascalCase or snake_case
* 🎉 feat: Min and Max function and options parameter added to use in error prints
