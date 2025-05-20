# Gocha
Live chat website for messageing in different chats
> Currenlty it is a personal project under development 
## features:
- **authentication**
- **login / signup**
- **create / delete chats**
- **join / leave chats**
- **pivate and public chats**
- **live messaging**
- **light / dark theme switching**
- **profile info**
## Application structure
**The Backend** server is built with golang and uses jwt authentication tokens, it has several packages like logging and validating and a database package using **Postgressql** for storing the user information and chats and messages and tokens and etc...

For the **Websocket connections** Gorilla/websocket package is used with managing and storing the connections in Memory for live messaging

**The Frontend** is built using **React/JS** that handles all the UI and the requests and connections to the backend. And uses many packages like react-router for managing different routes and authorization, and react-redux for storing state like user and chats and messages, and react-icons for UI icons.
## Demo
[![Gocha demo](https://img.youtube.com/vi/PpjK_zWgtbM/0.jpg)](https://www.youtube.com/watch?v=PpjK_zWgtbM)

