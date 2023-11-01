# iOS Challenge

This challenge is part of the iOS hiring process at [Heart Hands,
Inc](https://hearthands.tech/).

## Why this challenge?

Heart Hands is operating with a small team of talented people. We are looking
for seasoned engineers with strong technical foundation, deep knowledge of their
technical stack, and good product intuitions, that enjoy working on consumer
apps.

This challenge has been designed to give a glimpse of what it might be like when
joining the team. And the kind of technical challenges we face and care about.
We are expecting you to spend no more than 48 hours on this.

## Summary

You are tasked to develop a consumer application that allows you to converse
with a chatbot.

A server is available, you can read more about it in [`./server`](./server). The
documentation contains informations on how to run the server and what kinds of
API endpoints are available.

## Requirements

- The application should be composed of a single screen that list messages
  (imagine a WhatsApp conversation)
- This screen should display the messages as returned by the server 
- The app should allow sending new messages
- The server will send you messages that you must display

## Bonus

- Make your app work offline (both for app state and sending)
- Make your app resilient to bad network conditions (retries & timeouts)
- Make your app idempotent for both what you send and what you receive
- Make the app compatible to run on iPad and macOS
- Add animations to make the app shiny
