import { useState,useMemo, useCallback } from "react";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { useDisclosure } from "@heroui/modal";
import { CheckboxGroup } from "@heroui/checkbox";
import { Button } from "@heroui/button";
import { Input } from "@heroui/input";
import toast from "react-hot-toast";
import TemplateCard from "@/components/card/template";
import AddTemplate from "@/components/modal/template";
import { DeleteTemplate } from "@/utils/configuration/template";
import { SetTemplate } from "@/utils/select";
export default function TemplateLayout(props: { templates: Array<{name: string; [key: string]: any;}>; setUpdate: (update: boolean) => void; token: string; theme: string;}) {
  const { templates, setUpdate, token, theme } = props;
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [edit, setEdit] = useState(false);
  const [edit_template, setEditTemplate] = useState<{name: string; [key: string]: any;}>({name: ""});
  const [selected, setSelected] = useState<Array<string>>([]);
  const [search, setSearch] = useState<string>("");
  const [display_templates, setDisplayTemplate] = useState<Array<{name: string; [key: string]: any;}>>([]);
  useMemo(() => {
          templates &&
          search &&
          setDisplayTemplate(
              templates.filter((template) => template.name.includes(search))
          );
          templates && search === "" && setDisplayTemplate(templates);
      }, [search, templates]);
  const deleteItem = useCallback((value: Array<string>) => toast.promise(DeleteTemplate(token, value), {
          loading: "正在删除所选机场...",
          success: (res) => {
              res !== false ? res.map(
                  item => item.status ? toast.success(item.message) : toast.error(item.message)) : 
                  toast.error("删除失败, 未知错误")
              setUpdate(true)
              setSelected([])
              return "删除操作完成";
          },
          error: (e) => { 
              setUpdate(true)
              setSelected([])
              return e.code === "ERR_NETWORK" ? "请检查网络连接" : 
                  e.response.data.message ? e.response.data.message : e.response.data},
  }), []);
  const setItem = useCallback((value: string) => toast.promise(SetTemplate(token, value), {
          loading: "正在删除所选机场...",
          success: (res) => res,
          error: (e) => e.code === "ERR_NETWORK" ? "请检查网络连接" : 
                  e.response.data.message ? e.response.data.message : e.response.data
        }), []);
  const editItem = useCallback((value: {name: string, [key: string]: any}, edit: boolean) => {
    onOpen()
    setEdit(edit);
    setEditTemplate(value);
  },[])
  return (
    <div className="h-full w-full p-2">
      <AddTemplate isOpen={isOpen} onClose={onClose} edit={edit} token={token} template={edit_template} theme={theme} setUpdate={setUpdate} />
      <header className="flex gap-2 flex-row items-center p-2">
        <Input
          size="sm"
          variant="underlined"
          className="w-56"
          type="text"
          label={<span className="font-black">模板名称</span>}
          isClearable
          value={search}
          onValueChange={setSearch}
        />
        <Button
          size="sm"
          color="primary"
          variant="shadow"
          onPress={() => editItem({name: "111"}, false)}
        >
          <span className="font-black text-xl">添加</span>
        </Button>
        <Button size="sm" color="danger" variant="shadow" onPress={() => deleteItem(selected)}>
          <span className="font-black text-xl">删除</span>
        </Button>
      </header>
      <CheckboxGroup value={selected} onValueChange={setSelected}>
        <ScrollShadow style={{height: `calc(100% - 3rem)`}} className="w-full flex flex-wrap gap-2">
          {display_templates && display_templates.map((template) => (
            <div className="w-72" key={template.name}>
              <TemplateCard template={template} theme={theme} token={token} setDelete={deleteItem} setTemplate={setItem} editTemplate={() => editItem(template, true)}/>
            </div>
          ))}
        </ScrollShadow>
      </CheckboxGroup>
  </div>
  )
}