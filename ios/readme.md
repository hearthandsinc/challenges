# iOS Challenge

This challenge is part of the iOS hiring process at [Heart
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

You are tasked with the implementation of an iOS messaging app that allows to
communicate (send and receive text messages) with bots, each in their own 1:1
chat.

A server is available for you to use. You can read more about it in
[`./server`](./server). Its documentation contains informations on how it can be
run, and what kinds of API endpoints & entities are available.

Functional requirements:

- [ ] The app should start on a screen showing the list of all chats
- [ ] The app should allow opening each chat individually
- [ ] The app should allow sending messages to a chat
- [ ] The app should reflect the messages sent to and received from the server
- [ ] The app should be resilient to bad network conditions (drops & timeouts)

## Bonus

Some topics to look at to dive deeper:

- [ ] Make the app work offline (both for app state and sending)
- [ ] Make the app idempotent in regards to what you send and receive
- [ ] Integrate a splashscreen to hide chats while the app is loading
- [ ] Add support for optimistic sending to give instantaneity in the UI
- [ ] Add support for a local read/unread indicator
- [ ] Avoid block changing states so the app feels fluid & snappy
- [ ] Make the app compatible to run on iPad and macOS
- [ ] Make the app runnable on multiple devices
- [ ] _Anything_ that you feel could improve the UX!

## Design

We have prepared a design to help you during this challenge. You will also
receive the Figma link along with the challenge instructions.

![design](./design.png)

## Challenge Review

We know you only have a limited time alloted to deliver this challenge, and thus
will have to prioritize what you work on. A few things that are important for us
and that will be considered during the review:
- collaboration: is the code easy to read, maintain, and evolve?
- features: what did you prioritize to maximize your impact?
- testability: is the code tested or easily testable?
- documentation: is the readme clear? are important parts of the code documented? can we follow your thought process by looking at the git history?

Good luck, and enjoy!
