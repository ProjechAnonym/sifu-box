import { useEffect, useState } from "react";
import { NavigateFunction, useNavigate } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "@/redux/hooks";
import {
  Button,
  Divider,
  Checkbox,
  CheckboxGroup,
  ScrollShadow,
  Badge,
  Tooltip,
  Modal,
  ModalContent,
  ModalBody,
  ModalHeader,
  ModalFooter,
  useDisclosure,
} from "@nextui-org/react";
import toast from "react-hot-toast";
import Load from "@/components/load";
import ModalAdd from "@/layout/modalAdd";
import { ClienAuth } from "@/utils/ClientAuth";
import { RefreshConfig } from "@/utils/proxy/RefreshConfig";
import { FetchProxy } from "@/utils/proxy/FetchProxy";
import { DeleteItems, Key2Index } from "@/utils/proxy/DeleteProxy";
import { AddFile, AddProxy } from "@/utils/proxy/AddItems";
import { ProviderValue, RulesetValue } from "@/types/proxy";
function redirectLogin(navigate: NavigateFunction) {
  navigate("/login");
  toast.error("Please login", { duration: 2000 });
}
export default function Proxy() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const status = useAppSelector((state) => state.auth.status);
  const login = useAppSelector((state) => state.auth.login);
  const load = useAppSelector((state) => state.auth.load);
  const secret = useAppSelector((state) => state.auth.secret);
  const dark = useAppSelector((state) => state.mode.dark);
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  const [errors, setErrors] = useState<Array<string> | null>(null);
  const [file, setFile] = useState<FileList | null>(null);
  const [renewProxy, setRenewProxy] = useState(true);
  const [providers, setProviders] = useState<Array<ProviderValue>>([]);
  const [selectProviders, setSelectProviders] = useState<Array<string>>([]);
  const [selectRulesets, setSelectRulesets] = useState<Array<string>>([]);
  const [rulesets, setRulesets] = useState<Array<RulesetValue>>([]);
  const [proxy, setProxy] = useState<{
    providers: Array<ProviderValue>;
    rulesets: Array<RulesetValue>;
  } | null>(null);
  useEffect(() => {
    !login && !status && dispatch(ClienAuth({ auto: true }));
    login && !status && !load && redirectLogin(navigate);
    errors && errors.map((error) => toast.error(error));
    if (renewProxy && status && secret !== "") {
      FetchProxy(secret)
        .then((proxy) => {
          setProviders(proxy.Providers);
          setRulesets(proxy.Rulesets);
        })
        .catch(() => toast.error("获取配置信息失败"));
      setRenewProxy(false);
    }
  }, [status, login, load, secret, errors, renewProxy]);
  return (
    <main className="h-full w-full flex flex-col gap-y-4 p-3">
      <Load show={load} fullscreen={true} />
      <Modal
        isOpen={isOpen}
        onOpenChange={onOpenChange}
        isDismissable={false}
        placement="center"
        backdrop="blur"
        size="xl"
      >
        <ModalContent
          className={`${
            dark ? "sifudark" : "sifulight"
          } bg-background text-foreground`}
        >
          {(onClose) => (
            <form
              onSubmit={(e) => {
                e.preventDefault();
                proxy &&
                (proxy.rulesets.length === 0 || proxy.providers.length === 0)
                  ? toast.promise(AddProxy(secret, proxy), {
                      loading: "loading",
                      success: () => {
                        setRenewProxy(true);
                        setProxy(null);
                        onClose();
                        return "添加配置成功";
                      },
                      error: (err) => {
                        setRenewProxy(true);
                        err.code === "ERR_NETWORK"
                          ? setErrors(["网络错误"])
                          : setErrors(err.response.data.message);
                        setProxy(null);
                        return "出错了,请检查网络";
                      },
                    })
                  : toast.error("请先添加配置");
              }}
            >
              <ModalHeader>添加链接和规则集</ModalHeader>
              <ModalBody>
                <ModalAdd onProxySubmit={(proxy) => setProxy(proxy)} />
              </ModalBody>
              <ModalFooter>
                <Button
                  color="danger"
                  variant="light"
                  onPress={onClose}
                  size="sm"
                >
                  <span className="text-lg font-black">Close</span>
                </Button>
                <Button color="primary" size="sm" type="submit">
                  <span className="text-lg font-black">Confirm</span>
                </Button>
              </ModalFooter>
            </form>
          )}
        </ModalContent>
      </Modal>
      <header className="p-0 flex flex-wrap gap-x-4 gap-y-2 h-fit items-center">
        <div className="flex flex-col gap-y-2">
          <div className="flex flex-row gap-x-1">
            <Button size="sm" color="primary" onPress={onOpen}>
              <span className="text-lg font-black">添加</span>
            </Button>
            <Tooltip content="根据配置重新生成json文件">
              <Button
                size="sm"
                color="primary"
                onPress={() =>
                  toast.promise(RefreshConfig(secret), {
                    loading: "loading",
                    success: () => {
                      setErrors(null);
                      return "更新配置操作完成";
                    },
                    error: (e) => {
                      e.code === "ERR_NETWORK"
                        ? setErrors(["网络错误"])
                        : setErrors(e.response.data.message);
                      return "更新配置失败";
                    },
                  })
                }
              >
                <span className="text-lg font-black">刷新</span>
              </Button>
            </Tooltip>
            <Tooltip content="删除选中的机场链接和规则集">
              <Button
                size="sm"
                color="danger"
                onPress={() => {
                  const { urlsIndex, rulesetsIndex } = Key2Index(
                    selectProviders,
                    selectRulesets
                  );
                  if (rulesetsIndex.length === 0 && urlsIndex.length === 0) {
                    toast.error("请选择要删除的机场链接或规则集");
                    return;
                  }
                  toast.promise(DeleteItems(secret, urlsIndex, rulesetsIndex), {
                    loading: "loading",
                    success: (res) => {
                      res && setRenewProxy(true);
                      return "删除成功";
                    },
                    error: (err) => {
                      setRenewProxy(true);
                      err.code === "ERR_NETWORK"
                        ? setErrors(["网络错误"])
                        : setErrors(err.response.data.message);
                      return "删除失败";
                    },
                  });
                }}
              >
                <span className="text-lg font-black">删除</span>
              </Button>
            </Tooltip>
            <form
              className="flex flex-row justify-center items-center gap-x-2"
              onSubmit={(e) => {
                e.preventDefault();
                file
                  ? toast.promise(AddFile(secret, file), {
                      loading: "loading",
                      success: (res) => {
                        res && setRenewProxy(true);
                        setFile(null);
                        return "成功添加yaml文件";
                      },
                      error: (err) => `${err.response.data.message}`,
                    })
                  : toast.error("请选择文件");
              }}
            >
              <Badge
                content={file?.length || 0}
                color="success"
                placement="bottom-right"
              >
                <label
                  htmlFor="file-upload"
                  className="px-2 bg-primary h-8 rounded-md"
                >
                  <i className="bi bi-filetype-yml text-2xl" />
                </label>
              </Badge>
              <input
                type="file"
                id="file-upload"
                className="hidden"
                onChange={(e) => setFile(e.target.files)}
                multiple
              />
              <Button
                type="submit"
                size="sm"
                startContent={<i className="bi bi-upload text-xl" />}
              >
                <span className="text-lg font-black">上传</span>
              </Button>
            </form>
          </div>
        </div>
      </header>
      <Divider />
      <div className="h-1/2 w-full">
        {providers ? (
          <ScrollShadow className="h-full px-2">
            <CheckboxGroup
              label={
                <span className="font-black text-foreground">机场链接列表</span>
              }
              color="primary"
              value={selectProviders}
              onValueChange={setSelectProviders}
              className="w-full"
              orientation="horizontal"
              classNames={{ wrapper: "gap-x-4 gap-y-4" }}
            >
              {providers.map((provider, i) => (
                <Checkbox
                  key={`${provider.name}-${i}`}
                  value={`${provider.name}-${provider.id}`}
                  size="lg"
                  className="p-1"
                >
                  <div
                    className={`${
                      dark ? "bg-zinc-700" : "bg-slate-200"
                    } px-2 rounded-md flex flex-row gap-2`}
                  >
                    <span className="font-black">{provider.name}</span>
                    <span className="space-x-2">
                      {provider.proxy && (
                        <i className="bi bi-send-fill text-sm" />
                      )}
                      {provider.remote ? (
                        <i className="bi bi-hdd-network-fill text-sm" />
                      ) : (
                        <i className="bi bi-hdd-fill text-sm" />
                      )}
                    </span>
                  </div>
                </Checkbox>
              ))}
            </CheckboxGroup>
          </ScrollShadow>
        ) : (
          <div className="w-full h-full flex justify-center items-center text-4xl font-black">
            没有机场信息
          </div>
        )}
      </div>
      <Divider />
      <div className="h-1/2 w-full">
        {rulesets.length !== 0 ? (
          <ScrollShadow className="h-full px-2">
            <CheckboxGroup
              label={
                <span className="font-black text-foreground">规则集列表</span>
              }
              color="primary"
              value={selectRulesets}
              onValueChange={setSelectRulesets}
              className="w-full"
              orientation="horizontal"
              classNames={{ wrapper: "gap-x-4 gap-y-4" }}
            >
              {rulesets.map((ruleset, i) => (
                <Checkbox
                  value={`${ruleset.tag}-${ruleset.id}`}
                  key={`${ruleset.tag}-${i}`}
                  size="lg"
                  className="p-1"
                >
                  <div
                    className={`${
                      dark ? "bg-zinc-700" : "bg-slate-200"
                    } px-2 rounded-md flex flex-row gap-2`}
                  >
                    <span className="font-black">{ruleset.tag}</span>
                    <span className="space-x-1">
                      {!ruleset.china && (
                        <i className="bi bi-send-fill text-sm" />
                      )}
                      {ruleset.type === "remote" ? (
                        <i className="bi bi-hdd-network-fill text-sm" />
                      ) : (
                        <i className="bi bi-hdd-fill text-sm" />
                      )}
                    </span>
                  </div>
                </Checkbox>
              ))}
            </CheckboxGroup>
          </ScrollShadow>
        ) : (
          <div className="w-full h-full flex justify-center items-center text-4xl font-black">
            没有规则集信息
          </div>
        )}
      </div>
    </main>
  );
}
