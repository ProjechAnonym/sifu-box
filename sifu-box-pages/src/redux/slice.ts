import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import type { RootState } from "@/redux/store";
import { ClienAuth } from "@/utils/ClientAuth";
export const authSlice = createSlice({
  name: "auth",
  // `createSlice` will infer the state type from the `initialState` argument
  initialState: {
    secret: "",
    status: false,
    load: false,
    err: "",
    login: false,
  },
  reducers: {
    // Use the PayloadAction type to declare the contents of `action.payload`
    setSecret: (state, action: PayloadAction<string>) => {
      state.secret = action.payload;
    },
    setStatus: (state, action: PayloadAction<boolean>) => {
      state.status = action.payload;
    },
    setLoad: (state, action: PayloadAction<boolean>) => {
      state.load = action.payload;
    },
    setErr: (state, action: PayloadAction<string>) => {
      state.err = action.payload;
    },
    setLogin: (state, action: PayloadAction<boolean>) => {
      state.login = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(ClienAuth.pending, (state) => {
        state.load = true;
        state.login = true;
      })
      .addCase(
        ClienAuth.fulfilled,
        (
          state,
          action: PayloadAction<{
            status: boolean;
            msg: string;
            secret?: string;
          }>
        ) => {
          state.load = false;
          state.status = action.payload.status;
          !action.payload.status && (state.err = "login failed!");
          !action.payload.status && localStorage.removeItem("secret");
          action.payload.status && (state.secret = action.payload.secret || "");
          action.payload.status &&
            localStorage.setItem("secret", state.secret || "");
        }
      )
      .addCase(ClienAuth.rejected, (state) => {
        state.load = false;
        state.status = false;
        state.err = "login failed!";
        localStorage.removeItem("secret");
      });
  },
});

export const darkSlice = createSlice({
  name: "mode",
  initialState: { dark: true },
  reducers: {
    setDark: (state) => {
      state.dark = !state.dark;
    },
  },
});

export const { setSecret, setErr, setLoad, setLogin, setStatus } =
  authSlice.actions;
export const { setDark } = darkSlice.actions;
// Other code such as selectors can use the imported `RootState` type
export const selectSecret = (state: RootState) => state.auth.secret;
export const selectStatus = (state: RootState) => state.auth.status;
export const selectErr = (state: RootState) => state.auth.err;
export const selectLoad = (state: RootState) => state.auth.load;
export const selectLogin = (state: RootState) => state.auth.login;
export const selectDark = (state: RootState) => state.mode.dark;
