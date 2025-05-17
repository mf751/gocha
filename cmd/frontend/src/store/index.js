import { configureStore } from "@reduxjs/toolkit";
import userReducer from "./slices/user";
import chatsReducer from "./slices/chats";

const store = configureStore({
  reducer: {
    user: userReducer,
    chats: chatsReducer,
  },
});

export default store;
