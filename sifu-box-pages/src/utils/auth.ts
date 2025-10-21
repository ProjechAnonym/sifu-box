import axios from "axios";
import { createAsyncThunk, AsyncThunk } from "@reduxjs/toolkit";
export const Login: AsyncThunk<
  { status: boolean; token: string; admin: boolean },
  { user: string; password: string; admin: boolean; code: string },
  {}
> = createAsyncThunk<
  { status: boolean; token: string; admin: boolean },
  { user: string; password: string; admin: boolean; code: string }
>("auth/login", async (props) => {
  const { password, user, admin, code } = props;
  const data = new FormData();
  if (admin) {
    data.append("username", user);
    data.append("password", password);
  } else {
    data.append("code", code);
  }
  try {
    const res = await axios.post(
      `http://192.168.10.6:9090/api/login/${admin ? "admin" : "visitor"}`,
      data
    );
    console.log(res.data);
    return res.status === 200
      ? {
          status: true,
          token: res.data.message.jwt,
          admin: res.data.message.admin,
        }
      : {
          status: false,
          token: res.data.message.token,
          admin: false,
        };
  } catch (e) {
    console.error(e);
    throw e;
  }
});

export const Verify: AsyncThunk<
  { status: boolean; token: string; admin: boolean },
  {},
  {}
> = createAsyncThunk<{ status: boolean; token: string; admin: boolean }, {}>(
  "auth/verify",
  async () => {
    const jwt= localStorage.getItem("jwt");

    if (jwt === null) {
      return { status: false, token: "", admin: false };
    }
    try {
      const res = await axios.get("http://192.168.10.6:9090/api/verify", {
        headers: { Authorization: jwt },
      });
      return res.status === 200
        ? {
            status: true,
            token: res.data.message.jwt,
            admin: res.data.message.admin,
          }
        : {
            status: false,
            token: res.data.message.jwt,
            admin: false,
          };
    } catch (e) {
      console.error(e);
      throw e;
    }
  }
);