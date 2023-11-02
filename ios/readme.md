# iOS Challenge

This challenge is part of the iOS hiring process at [Heart
Hands](https://hearthands.tech/).

## Why this challenge?

Heart Hands is operating with a small team of dedicated and talented people. We
are looking for seasoned engineers with a deep technical knowledge, strong
understanding of their technical stack, and good product intuitions to join our
team.

This challenge has been designed to give a glimpse of what it is like to build a
messaging app, and the kind of technical challenges we face and care about. We
are expecting you to spend no more than 48 hours on this.

## Summary

You are tasked to develop a messaging app that allows to send and receive
messages with several bots, each in their own chat.

A server is available for you to use, you can read more about it in
[`./server`](./server). The documentation contains informations on how it can be
run and what kinds of API endpoints are available.

## Requirements

This is the base foundation under which we consider the app usable:

- [ ] The app should start on a screen showing the list of all chats
- [ ] The app should allow opening each chat individually
- [ ] The app should allow sending messages to a chat
- [ ] The app should reflect the latest messages sent to and received from the server

## Bonus

Some ideas worth exploring to improve the experience:

- [ ] Make the app work offline (both for app state and sending)
- [ ] Make the app resilient to bad network conditions (retries & timeouts)
- [ ] Make the app idempotent in regards to what you send and receive
- [ ] Add support for optimisic sending to give instantaneity in the UI
- [ ] Add support for a local read/unread indicator
- [ ] Avoid block changing states so the app feels fluid & snappy
- [ ] Make the app compatible to run on iPad and macOS
- [ ] Make the app runnable on multiple devices

## Design

We have prepared a design to help you when working on this challenge.

[![design](./design.png)](https://www.figma.com/file/jyaVdqka2yAZ53b5YzkN6T/iOS-Challenge)
