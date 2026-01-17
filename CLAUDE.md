# Project Context
GoGolf is a simple command line based golf game in the RPG genre.

The original inspiration for this game was the table top rpg experience of RuneScape: Kingdoms the roleplaying game. 
Because of this inspiration, the basic mechanics of the game are as follows:

* Dice roll mechanics to determine success/fail of a player action
    * 3d6 dice rolls
    * success or fail is determined by checking the margin of the roll, modifiers, and the target number
    * target number is determined by the skills involved in the check as well as the difficulty of the action
    * success and fail is not binary, it has tiers of success/fail including 'critical' when all 3 dice match, or the margin is significant enough
* Character building through 'paper doll'
    * A character has skills and abilities, such as strength, accuracy, green reading, irons, drivers, etc...
    * Both equipment and skills/abilities are used to determine the outcome of actions
    * Equipment, such as clubs, balls, gloves, can be purchased from the ProShop
    * Skills/Abilities gain experience when used for a skill check
    * Skills/Abilities can be leveled to a maximum level of 9
* The course is a grid of locations the ball can be within
    * Each location will have properties that determine the lie of the ball and will influence the difficulty of the shot
        * ex: a location may be considered to be 'in a bunker', 'in the rough', 'or in the fairway', providing different challenges and outcomes

# Development Best Practices
* Always use TDD (test driven development).
    * Write tests based on expected input/output pairs. Avoid creating mock implementations, even for functionality that doesn’t exist yet in the codebase
    * Verify the test fails
    * Commit failing test(s)
    * Implement just enough code to pass test(s)
    * Commit working code