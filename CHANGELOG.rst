CHANGELOG
=========

UNRELEASED
----------

2.1.2 (2025-07-13)

* 🐛 fix: required rule doesn't check for zero values

2.1.1 (2025-07-11)

* 🐛 fix: password validation was a little bit limited

2.1.0 (2025-06-07)

* 🎉 feat: added GetRule, GetStructRule in ruleSet and GetValidator with GetChildrenValidator in Validator

2.0.0 (2024-05-07)
------------------

* 🎉 feat: a context goes through the validators now

1.4.4 (2024-05-07)
------------------

* 🐛 fix: fixed #11 reported bug

1.4.3 (2023-10-10)
------------------

* 🐛 fix: add custom messages for rules with params is now possible

1.4.2 (2023-08-21)
------------------

* 🐛 fix: unexported fields will get ignored to prevent unwanted behavior

1.4.1 (2023-08-20)
------------------

* 🐛 fix: WhenNotExistAll and WhenNotExistOne were not usable in tags + Validator does not add a type check validator anymore + bug fixed on after passing required

1.4.0 (2023-08-20)
------------------

* ✅ test: test cases updated for new features
* 🎉 feat: field names will get translated
* 🎉 feat: ability to pass custom validators from generators to ruleSet is now possible
* 🎉 feat: added AlwaysCheckRules method in ruleSet
* 🎉 feat: added WhenNotExistOne + WhenNotExistAll

1.3.5 (2023-08-19)
------------------

* 🐛 fix: bug fixed on setting error messages for conflicted keys in tags

1.3.4 (2023-07-11)
------------------

* 🐛 fix: children rules can now have no value (they had before for some reason)

1.3.3 (2023-07-10)
------------------

* 🐛 fix: fixed problem on Phone number rule + fixed bug on setDefaultOn function

1.3.2 (2023-05-02)
------------------

* 🐛 fix: phone number rule will be just valid with international phone numbers formats

1.3.1 (2023-04-30)
------------------

* 🐛 fix: mixed a bug in phonenumber verification rule and made it better

1.3.0 (2023-01-28)
------------------

* 🎉 feat: added new functions to use in patch method or other scenarios:
1. `SetDefaultOnNil` is added which is used in PATCH method.
2. `SetDefault` is added which is used in scenarios that you want to have default passed values in nil, zero or empty of element positions.

1.2.0 (2022-12-20)
------------------

* 🎉 feat: translator feature added

1.1.2 (2022-12-16)
------------------

* 🐛 fix: some odd behaviors of creating a validator from a generator fixed
* 📖 docs: documentation updated

1.1.1 (2022-12-11)
------------------

* ✅ test: added test cases for ComplexValidator method
* 🎉 feat: ComplexValidator added to Generator for creating a complex(struct or map) from zero level

1.1.0 (2022-11-29)
------------------

* ✅ test: test cases for new feature added
* 🎉 feat: support for errors returned from [gin](https://github.com/gin-gonic/gin)'s Bind actions added

1.0.0 (2022-11-09)
------------------

* 🎉 feat: added CustomMessages method for generator
* 🎉 feat: #1 feature added, added `c.` operator in tags to use when need to define rules for children of a slice element.
* 🐛 fix: #4

0.0.2 (2022-11-05)
------------------

* 🐛 fix: #1, Pointers of different variables when used inside `g.Validator` or `Validator.Validate` methods, didn't work properly.

0.0.1 (2022-10-30)
------------------

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
