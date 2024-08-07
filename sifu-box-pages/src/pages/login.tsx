import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import { Input, Button } from "@nextui-org/react";
import Load from "@/components/load";
import { ClienAuth } from "@/utils/ClientAuth";
import { setErr } from "@/redux/slice";
export default function Login() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const load = useAppSelector((state) => state.auth.load);
  const status = useAppSelector((state) => state.auth.status);
  const err = useAppSelector((state) => state.auth.err);
  const [password, setPassword] = useState("");
  const [invalid, setInvalid] = useState(false);
  const [visible, setVisible] = useState(false);
  useEffect(() => {
    status && navigate("/");
  }, [status]);
  return (
    <form
      className="w-full h-full content-center"
      onSubmit={(e) => {
        e.preventDefault();
        dispatch(ClienAuth({ password: password, auto: false }));
      }}
    >
      <Load show={load} fullscreen={true} />
      <Input
        label="Password"
        type={visible ? "text" : "password"}
        size="lg"
        required={true}
        fullWidth={false}
        className="w-80 mx-auto"
        value={password}
        isInvalid={invalid || err !== ""}
        onValueChange={(e) => {
          e !== "" && setInvalid(false);
          err !== "" && dispatch(setErr(""));
          setPassword(e);
        }}
        errorMessage={err !== "" ? "login failed!" : "Please enter password"}
        endContent={
          <button
            className="focus:outline-none"
            type="button"
            onClick={() => setVisible(!visible)}
          >
            {visible ? (
              <i className="bi bi-eye-slash-fill" />
            ) : (
              <i className="bi bi-eye-fill" />
            )}
          </button>
        }
      />
      <p className="h-fit w-80 my-4 gap-4 flex flex-row justify-end mx-auto">
        <Button type="submit" color="primary" onPress={() => setInvalid(true)}>
          <span className="font-black text-xl">确认</span>
        </Button>
        <Button
          type="button"
          color="danger"
          onPress={() => {
            setPassword("");
            setInvalid(false);
            dispatch(setErr(""));
          }}
        >
          <span className="font-black text-xl">清空</span>
        </Button>
      </p>
    </form>
  );
}
