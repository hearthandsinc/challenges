# AI Challenge

This challenge is part of the AI hiring process at Heart Hands.

## Why this challenge?

Heart Hands is operating with a small team of dedicated & talented people. We
are looking for seasoned engineers with deep technical knowledge, a strong
understanding of their technical stack, and excellent product intuitions to join
our team.

This exercise has been designed to give a glimpse of what it is like to build an
AI assistant and the kind of technical challenges we face and care about. We are
expecting you to spend between 4 and 6 hours on this challenge.

## Instructions

You are tasked with implementing a retrieval system for a messaging app. It must
expose API endpoints that enable users to query their conversation history using
natural language. You are free to implement the API in any way you feel the most
comfortable with. You can find a dataset of chat conversations between several
users
[here](https://github.com/hearthandsinc/challenges/blob/main/ai/chats.json).

We enforce no technical constraints: you are free to choose the language, data
layer, network protocol, and design your API as you see fit. You are
purposefully being given a lot of freedom here, and you will not be judged on
these decisions alone, but we will challenge the understanding of the trade-offs
you make.

You can find a simple boilerplate for a Go/PGSQL/Qdrant app
[here](https://github.com/hearthandsinc/challenges/tree/main/ai/boilterplate),
but again, you are free to use any technology you want.

Functional requirements:

- [ ] Users can search through their chat history using natural language. To
  evaluate the solution, we'll be asking questions from Matthias, Aymeric and
  David's point of view:
- Matthias' id is: `018f2685-a936-44a9-8ff4-9ef0c98289b8`
- Aymeric's id is: `82247543-5c2d-46d9-9a0d-fe25482922b5`
- David's id is: `2d614bef-2b01-4021-b63a-9d04658536f3`
- Example of the questions we will be asking:
    - "What did Aymeric and I discuss last week?"
    - "Where does Aymeric live?"
    - "What did David tell me about the roadmap?"
- [ ] Your solution is multitenant. Users can only search their own history using their user ID.
- [ ] Users can query their chat history through an API endpoint.

## Bonus

Some topics that we find interesting to dig into:

- [ ] Ability to target a specific date.
E.g.:
    - "Can you recap last Thursday's chat activity?"
    - "Can you summarize our discussion with David from March 12th?"
- [ ] Ability to count occurrences.
E.g.:
    - "How many messages did Aymeric and I exchange so far?"
- [ ] The assistant chat keeps track of one request to the next.
E.g.:
    1. "What are Aymeric's requirements for buying a new home?"
    2. "What are David's?"

## Challenge Review

We know you only have a limited time allotted to deliver this challenge, and thus will have to prioritize what you work on. A few things that are important for us and that will be considered during the review:

- **documentation**: is the README clear? are important parts of the code documented?
- **impact**: which features did you prioritize?
- **maintainability**: is the code well-structured and easy to read/evolve?
- **robustness**: are edge cases considered?

Good luck, and enjoy!
