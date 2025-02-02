import { configureStore } from "@reduxjs/toolkit";
import { themeSlice, authSlice } from "./slice";
// ...

export const store = configureStore({
  reducer: {
    theme: themeSlice.reducer,
    auth: authSlice.reducer,
  },
});

// 从 store 本身推断 `RootState` 和 `AppDispatch` 类型
export type RootState = ReturnType<typeof store.getState>;
// 推断类型：{posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch;
