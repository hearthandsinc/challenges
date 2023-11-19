# Backend Challenge

This challenge is part of the backend hiring process at [Heart
Hands](https://hearthands.tech/).

## Why this challenge?

Heart Hands is operating with a small team of dedicated & talented people. We
are looking for seasoned engineers with a deep technical knowledge, strong
understanding of their technical stack, and excellent product intuitions to join
our team.

This exercise has been designed to give a glimpse of what it is like to build a
messaging app, and the kind of technical challenges we face and care about. We
are expecting you to spend between 4 and 6 hours on this challenge.

## Instructions

You are tasked with the server-side implementation of a messaging app that
allows clients to send and receive text messages in private 1:1 chats.

We enforce no technical constraints: you are free to choose the language, data
layer, network protocol, and design your API as you see fit. You are
purposefully being given a lot of freedom here, and you will not be judged on
these decisions alone, but we will challenge the understanding of the trade-offs
you make.

Functional requirements:

- [ ] Clients should be able to uniquely identify themselves with a unique phone number (without authentication)
- [ ] Clients should be able to send messages to other clients
- [ ] Clients should be able to list their chats
- [ ] Clients should be able to list all messages in a chat

## Bonus

Some topics that we find interesting to dig:

- [ ] Implement SMS authentication to verify the clients phone numbers & add
      clients sessions
- [ ] Choose a network protocol that enables soft real-time message delivery to clients
- [ ] Support SMS forwarding to relay the messages to the clients phone numbers
- [ ] Make the server message ingestion idempotent
- [ ] Add support for sent/delivered/read message status
- [ ] Add support for chats and messages pagination

## Challenge Review

We know you only have a limited time alloted to deliver this challenge, and thus
will have to prioritize what you work on. A few things that are important for us
and that will be considered during the review:
- **impact**: which features did you prioritize?
- **maintainability**: is the code well-structured and easy to read/update?
- **testability**: is the code tested or easily testable?
- **documentation**: is the readme clear? are important parts of the code
  documented?

Good luck, and enjoy!
