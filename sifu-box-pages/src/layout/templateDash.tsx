import { useState, useMemo } from "react";
import {
  Input,
  Button,
  Checkbox,
  CheckboxGroup,
  Card,
  CardHeader,
  CardBody,
  CardFooter,
  useDisclosure,
} from "@nextui-org/react";
import SetTemplate from "@/components/setting/setTemplate";
import toast from "react-hot-toast";
import { DeleteTemplate } from "@/utils/template/DeleteTemplate";

export default function TemplateDash(props: {
  dark: boolean;
  secret: string;
  headHeight: number;
  template: Array<{
    Name: string;
    Template: Object;
  }> | null;
  baseTemplate: { Name: string; Template: Object } | null;
  setUpdateTemplate: (updateTemplate: boolean) => void;
}) {
  const {
    dark,
    secret,
    headHeight,
    template,
    setUpdateTemplate,
    baseTemplate,
  } = props;
  const { isOpen, onClose, onOpen } = useDisclosure();
  const [search, setSearch] = useState("");
  const [newTemplate, setNewTemplate] = useState(false);
  const [displayTemplate, setDisplayTemplate] = useState<Array<{
    Name: string;
    Template: Object;
  }> | null>(null);
  const [editTemplate, setEditTemplate] = useState<{
    Name: string;
    Template: Object;
  } | null>();
  const [selected, setSelected] = useState<Array<string>>([]);
  useMemo(() => {
    template &&
      search &&
      setDisplayTemplate(
        template.filter((value) => value.Name.includes(search))
      );
    template && !search && setDisplayTemplate(template);
  }, [search, template]);
  return (
    <div
      className="p-2 overflow-auto"
      style={{ height: `calc(100% - ${headHeight}px)` }}
    >
      <SetTemplate
        isOpen={isOpen}
        onClose={onClose}
        dark={dark}
        template={editTemplate ? editTemplate : null}
        newTemplate={newTemplate}
        secret={secret}
        setUpdateTemplate={setUpdateTemplate}
      />
      <header className="flex flex-row items-center gap-2 h-12">
        <Input
          variant="underlined"
          className="w-56"
          type="text"
          label={<span className="font-black">模板名称</span>}
          isClearable
          value={search}
          onValueChange={setSearch}
        />
        <Button
          color="primary"
          size="sm"
          onPress={() => {
            onOpen();
            setNewTemplate(true);
            baseTemplate && setEditTemplate(baseTemplate);
          }}
        >
          <span className="text-lg font-black">添加</span>
        </Button>
        <Button
          color="danger"
          size="sm"
          onPress={() =>
            selected.length &&
            toast.promise(DeleteTemplate(secret, selected), {
              loading: "loading",
              success: (res) => {
                setUpdateTemplate(true);
                return res
                  ? `${selected.join(",")}删除成功`
                  : `${selected.join(",")}删除失败`;
              },
              error: (err) =>
                err.code === "ERR_NETWORK"
                  ? "网络错误"
                  : err.response.data.message,
            })
          }
        >
          <span className="text-lg font-black">删除</span>
        </Button>
      </header>
      {displayTemplate && (
        <CheckboxGroup
          value={selected}
          onValueChange={setSelected}
          className="my-4"
        >
          <div className="flex flex-wrap gap-2">
            {displayTemplate.map((value, i) => (
              <Card
                className={`hover:cursor-pointer w-80`}
                key={`${value.Name}-${i}`}
                isHoverable
                shadow="none"
              >
                <CardHeader
                  className={`flex justify-between ${
                    dark ? "bg-zinc-800" : "bg-slate-100"
                  }`}
                >
                  <span className="font-black text-lg">{value.Name}</span>
                  {value.Name !== "default" && <Checkbox value={value.Name} />}
                </CardHeader>
                <CardBody
                  className={`flex flex-col gap-2 ${
                    dark ? "bg-zinc-800" : "bg-slate-100"
                  }`}
                >
                  <div className="h-48 overflow-auto whitespace-pre-wrap w-fit">
                    {JSON.stringify(value.Template, null, 4)}
                  </div>
                </CardBody>
                <CardFooter
                  className={`${
                    dark ? "bg-zinc-800" : "bg-slate-100"
                  } flex justify-end gap-2`}
                >
                  {value.Name !== "default" && (
                    <Button
                      size="sm"
                      color="danger"
                      onPress={() =>
                        toast.promise(DeleteTemplate(secret, [value.Name]), {
                          loading: "loading",
                          success: (res) => {
                            setUpdateTemplate(true);
                            return res
                              ? `${value.Name}删除成功`
                              : `${value.Name}删除失败`;
                          },
                          error: (err) =>
                            err.code === "ERR_NETWORK"
                              ? "网络错误"
                              : err.response.data.message,
                        })
                      }
                    >
                      <span className="font-black text-lg">删除</span>
                    </Button>
                  )}
                  <Button
                    color="primary"
                    size="sm"
                    onPress={() => {
                      onOpen();
                      setNewTemplate(false);
                      setEditTemplate(displayTemplate[i]);
                    }}
                  >
                    <span className="font-black text-lg">修改</span>
                  </Button>
                </CardFooter>
              </Card>
            ))}
          </div>
        </CheckboxGroup>
      )}
    </div>
  );
}
