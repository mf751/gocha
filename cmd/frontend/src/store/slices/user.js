import { createSlice } from "@reduxjs/toolkit";

const userSlice = createSlice({
  name: "user",
  initialState: { user: {}, loggedIn: false },
  reducers: {
    setUser(state, action) {
      state.user = action.payload;
    },
    setLoggedIn(state, action) {
      state.loggedIn = action.payload;
    },
  },
});

export const { setUser, setLoggedIn } = userSlice.actions;
export default userSlice.reducer;
