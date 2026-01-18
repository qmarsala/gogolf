# Project Context
GoGolf is a simple command line based golf game in the RPG genre.
Refer to [CORE_MECHANICS.md](./CORE_MECHANICS.md) to understand how the game works when considering new features.
Refer to [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md) to understand the current priorities.


# Development Best Practices
* Always start by making a PLAN.
    * The PLAN should clearly state the problem, potential solutions, and expected results
* Always use TDD (test driven development).
    * Write tests based on expected input/output pairs. Avoid creating mock implementations, even for functionality that doesn’t exist yet in the codebase
    * Verify the test fails
    * Commit failing test(s)
    * Implement just enough code to pass test(s)
    * Commit working code
* Keep comments to a minimum
    * Remove unhelpful comments when encountered
    * Prefer self-documenting code, with no comments
* Utilize Trunk-Based branching. Creating a branch and PR only when large changes need careful review