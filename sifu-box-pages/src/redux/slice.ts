import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import type { RootState } from "@/redux/store";
import { Login, Verify} from "@/utils/auth";
export const authSlice = createSlice({
  name: "auth",
  initialState: {
    admin: false,
    load: false,
    status: false,
    jwt: "",
    err: "",
    auto: true,
  },
  reducers: {
    setLoad: (state, action: PayloadAction<boolean>) => {
      state.load = action.payload;
    },
    setStatus: (state, action: PayloadAction<boolean>) => {
      state.status = action.payload;
    },
    setJwt: (state, action: PayloadAction<string>) => {
      state.jwt = action.payload;
    },
    setErr: (state, action: PayloadAction<string>) => {
      state.err = action.payload;
    },
    setAuto: (state, action: PayloadAction<boolean>) => {
      state.auto = action.payload;
    },
    setAdmin: (state, action: PayloadAction<boolean>) => {
      state.admin = action.payload;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(Login.pending, (state) => {
        state.load = true;
        state.auto = false;
      })
      .addCase(
        Login.fulfilled,
        (
          state,
          action: PayloadAction<{
            status: boolean;
            token: string;
            admin: boolean;
          }>
        ) => {
          state.load = false;
          state.status = action.payload.status;
          !action.payload.status &&
            (state.err = "登录失败, 请检查密码或者用户名");
          !action.payload.status && localStorage.removeItem("jwt");
          action.payload.status && (state.jwt = action.payload.token || "");
          action.payload.status &&
            localStorage.setItem("jwt", state.jwt || "");
          action.payload.status &&
            (state.admin = action.payload.admin || false);
        }
      )
      .addCase(Login.rejected, (state) => {
        state.load = false;
        state.status = false;
        state.err = "登录失败, 请检查密码或者用户名";
        localStorage.removeItem("jwtToken");
      })
      .addCase(Verify.pending, (state) => {
        state.load = true;
        state.auto = false;
      })
      .addCase(
        Verify.fulfilled,
        (
          state,
          action: PayloadAction<{
            status: boolean;
            token: string;
            admin: boolean;
          }>
        ) => {
          state.load = false;
          state.status = action.payload.status;
          !action.payload.status &&
            (state.err = "登录失败, 请检查密码或者用户名");
          !action.payload.status && localStorage.removeItem("jwtToken");
          action.payload.status && (state.jwt = action.payload.token || "");
          action.payload.status &&
            localStorage.setItem("jwtToken", state.jwt || "");
          action.payload.status &&
            (state.admin = action.payload.admin || false);
        }
      )
      .addCase(Verify.rejected, (state) => {
        state.load = false;
        state.status = false;
        state.err = "自动登录失败";
        localStorage.removeItem("jwtToken");
      });
  },
});
export const themeSlice = createSlice({
  name: "theme",
  initialState: { theme: "sifulight" },
  reducers: {
    setTheme: (state, action: PayloadAction<string>) => {
      state.theme = action.payload;
    },
  },
});
export const { setTheme } = themeSlice.actions;
export const { setLoad, setStatus, setJwt, setErr, setAuto, setAdmin } =
  authSlice.actions;
export const selectTheme = (state: RootState) => state.theme.theme;
export const selectLoad = (state: RootState) => state.auth.load;
export const selectStatus = (state: RootState) => state.auth.status;
export const selectJwt = (state: RootState) => state.auth.jwt;
export const selectErr = (state: RootState) => state.auth.err;
export const selectAuto = (state: RootState) => state.auth.auto;