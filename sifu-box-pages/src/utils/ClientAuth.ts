import axios from "axios";
import { createAsyncThunk, AsyncThunk } from "@reduxjs/toolkit";

export const ClienAuth: AsyncThunk<
  { status: boolean; msg: string; secret?: string },
  { password?: string; auto: boolean },
  {}
> = createAsyncThunk<
  { status: boolean; msg: string; secret?: string },
  { password?: string; auto: boolean }
>("auth/verify", async (props) => {
  const { password = "", auto } = props;
  const secret = localStorage.getItem("secret") || "";
  if (auto && secret === "") {
    return { status: false, msg: "failed" };
  }
  try {
    const res = await axios.get("/api/verify", {
      headers: { Authorization: auto ? secret : password },
    });
    return res.status === 200
      ? { status: true, msg: "success", secret: auto ? secret : password }
      : { status: false, msg: "failed" };
  } catch (e) {
    throw e;
  }
});
