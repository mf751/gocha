import { createSlice } from "@reduxjs/toolkit";

const chatsSlice = createSlice({
  name: "chats",
  initialState: { chats: [], loaded: false },
  reducers: {
    setChats(state = initialState, action) {
      state.chats = action.payload;
    },
    setLoaded(state = initialState, action) {
      state.loaded = action.payload;
    },
  },
});

export const { setChats, setLoaded } = chatsSlice.actions;
export default chatsSlice.reducer;
