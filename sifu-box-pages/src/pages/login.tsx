import { useState, useEffect } from "react";
import { useAppSelector, useAppDispatch } from "@/redux/hooks";
import { useNavigate } from "react-router-dom";
import { setErr } from "@/redux/slice";
import { Form } from "@heroui/form";
import { Card, CardBody, CardFooter, CardHeader } from "@heroui/card";
import { Input } from "@heroui/input";
import { Button } from "@heroui/button";

import { Login, AutoLogin } from "@/utils/auth";
import toast from "react-hot-toast";
import DefaultLayout from "@/layouts/default";
import Load from "@/components/load";
import { Switch } from "@heroui/switch";
export default function LoginPage() {
  const load = useAppSelector((state) => state.auth.load);
  const status = useAppSelector((state) => state.auth.status);
  const auto = useAppSelector((state) => state.auth.auto);
  const err = useAppSelector((state) => state.auth.err);
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const [user, setUser] = useState("");
  const [password, setPassword] = useState("");
  const [admin, setAdmin] = useState(true);
  const [code, setCode] = useState("");
  const [login, setLogin] = useState(false);
  useEffect(() => {
    status && navigate("/home");
    login && err !== "" && toast.error(err);
    err !== "" && dispatch(setErr(""));
    auto && dispatch(AutoLogin({}));
    !login && setLogin(true);
  }, [status, err, auto, login]);
  return (
    <DefaultLayout>
      <Load
        fullscreen={true}
        show={load}
        color="primary"
        label="loging..."
        label_color="primary"
      />
      <Form
        id="loginForm"
        validationBehavior="native"
        className="h-full w-full"
        onSubmit={(e) => {
          e.preventDefault();
          dispatch(
            Login({ user: user, password: password, code: code, admin: admin })
          );
        }}
        onReset={() => {
          setUser("");
          setPassword("");
        }}
      >
        <Card className="m-auto w-64 h-64">
          <CardHeader>
            <span className="text-xl font-black m-2 select-none">sifu-box</span>
          </CardHeader>
          <CardBody className="flex flex-col gap-4 items-center h-full">
            {admin && (
              <Input
                value={user}
                onValueChange={setUser}
                label="user"
                isRequired
                size="sm"
                isClearable
                errorMessage="用户名不能为空"
              />
            )}
            {admin && (
              <Input
                value={password}
                onValueChange={setPassword}
                label="password"
                type="password"
                isRequired
                size="sm"
                isClearable
                errorMessage="密码不能为空"
              />
            )}
            {!admin && (
              <Input
                value={code}
                onValueChange={setCode}
                label="code"
                type="code"
                isRequired
                size="sm"
                isClearable
                errorMessage="密钥不能为空"
              />
            )}
          </CardBody>
          <CardFooter className="flex flex-row justify-between">
            <Switch isSelected={admin} onValueChange={setAdmin} size="sm">
              <span className="font-black text-md">
                {admin ? "管理员" : "访客"}
              </span>
            </Switch>
            {/* <Link href="/reset" isBlock>
              忘记密码
            </Link> */}
            <div className="flex flex-row gap-2">
              <Button size="sm" variant="shadow" type="reset" color="danger">
                <span className="text-lg font-black">清空</span>
              </Button>
              <Button size="sm" variant="shadow" type="submit" color="primary">
                <span className="text-lg font-black">登录</span>
              </Button>
            </div>
          </CardFooter>
        </Card>
      </Form>
    </DefaultLayout>
  );
}
