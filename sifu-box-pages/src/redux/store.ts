import { configureStore } from "@reduxjs/toolkit";
import { authSlice } from "./slice";
import { darkSlice } from "./slice";
export const store = configureStore({
  reducer: {
    auth: authSlice.reducer,
    mode: darkSlice.reducer,
  },
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch;
