import { createSlice } from "@reduxjs/toolkit";

const chatsSlice = createSlice({
  name: "chats",
  initialState: { chats: [] },
  reducers: {
    setChats(state, action) {
      state.chats = action.payload;
    },
  },
});

export const { setChats } = chatsSlice.actions;
export default chatsSlice.reducer;
