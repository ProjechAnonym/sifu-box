import { useState, useMemo, useRef } from "react";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Badge } from "@heroui/badge";
import { Button } from "@heroui/button";
import { Checkbox, CheckboxGroup } from "@heroui/checkbox";
import { Input } from "@heroui/input";
import { useDisclosure } from "@heroui/modal";
import toast from "react-hot-toast";
import AddProviders from "@/components/card/provider";
import { DeleteProvider, AddProviderFiles } from "@/utils/configuration/provider";
import { Provider } from "@/types/setting/provider";
export default function ProviderLayout(props: {providers: Array<Provider>, theme: string, token: string, setUpdate: (update: boolean) => void}) {
    const { providers, theme, token, setUpdate } = props;
    const { isOpen, onOpen, onClose } = useDisclosure();
    const file_input = useRef<HTMLInputElement>(null);
    const [filtered_providers, setFilteredProviders] = useState(providers);
    const [checked, setChecked] = useState<Array<string>>([]);
    const [search, setSearch] = useState("");
    const [edit, setEdit] = useState(false);
    const [edit_value, setEditValue] = useState<{name: string, path: string, remote: boolean}>({name: "", path: "", remote: false});
    const [files, setFiles] = useState<FileList | null>(null);
    useMemo(() => {
        providers &&
        search &&
        setFilteredProviders(
            providers.filter((provider) => provider.name.includes(search))
        );
        providers && search === "" && setFilteredProviders(providers);
    }, [search, providers]);
    const deleteItem = (value: Array<string>) => toast.promise(DeleteProvider(token, value), {
        loading: "正在删除所选机场...",
        success: (res) => {
            res !== false ? res.map(
                item => item.status ? toast.success(item.message) : toast.error(item.message)) : 
                toast.error("删除失败, 未知错误")
            setUpdate(true)
            setChecked([])
            return "删除操作完成";
        },
        error: (e) => { 
            setUpdate(true)
            setChecked([])
            return e.code === "ERR_NETWORK" ? "请检查网络连接" : 
                e.response.data.message ? e.response.data.message : e.response.data},
    });
    const addFiles = (files: FileList) => toast.promise(AddProviderFiles(token, files), {
        loading: "正在添加所选机场...",
        success: (res) => {
            res ? res.map(item => item.status ? toast.success(item.message) : toast.error(item.message)) : toast.error("添加失败, 未知错误")
            setFiles(null)
            setUpdate(true)
            file_input.current ? file_input.current.value = "" : console.error("文件选择input元素出现未知错误")
            return "添加操作完成";
        },
        error: (e) => {
            setFiles(null)
            setUpdate(true)
            file_input.current ? file_input.current.value = "" : console.error("文件选择input元素出现未知错误")
            return e.code === "ERR_NETWORK" ? "请检查网络连接" : 
                e.response.data.message ? e.response.data.message : e.response.data}
    })
    return (
        <ScrollShadow className="w-full h-1/2">
            <AddProviders edit={edit} isOpen={isOpen} onClose={onClose} theme={theme} token={token} setUpdate={setUpdate} initial_value={edit_value} />
            <header className="flex flex-wrap gap-1 p-1 items-center">
                <Input
                    size="sm"
                    label={<span className="text-md font-black">机场</span>}
                    variant="underlined"
                    value={search}
                    onValueChange={setSearch}
                    className="w-24"
                />
                <Button color="danger" size="sm" variant="shadow" onPress={() => deleteItem(checked)}>
                    <span className="text-xl font-black">删除</span>
                </Button>
                <Button color="primary" size="sm" variant="shadow" onPress={()=>{setEdit(false); setEditValue({name: "", path: "", remote: false});onOpen()}}>
                    <span className="text-xl font-black">添加</span>
                </Button>
                <form className="flex flex-row gap-1" 
                    onSubmit={e => {
                        e.preventDefault()
                        files ? addFiles(files) : toast.error("请选择文件")
                    }}
                >
                    <Badge content={files ? files.length : 0} color="success" placement="bottom-right">
                        <label htmlFor="file-upload" className="px-2 bg-primary h-8 rounded-md hover:cursor-pointer hover:bg-opacity-85 transition-all text-white">
                            <i className="bi bi-filetype-yml text-2xl" />
                        </label>
                    </Badge>
                    <input
                        type="file"
                        id="file-upload"
                        className="hidden"
                        onChange={(e) => setFiles(e.target.files)}
                        multiple
                        ref={file_input}
                    />
                    <Button size="sm" color="primary" variant="shadow" type="submit">
                        <span className="font-black text-xl">上传</span>
                    </Button>
                </form>
            </header>
            <CheckboxGroup value={checked} onValueChange={setChecked}>
                <div className={`p-2 flex flex-wrap gap-2`}>
                    {filtered_providers && filtered_providers.map((provider) => (
                        <div key={provider.name} className="flex flex-row justify-center h-fit gap-1">
                            <Badge
                                content={<i className="bi bi-trash-fill" />}
                                placement="bottom-left"
                                color="danger"
                                shape="rectangle"
                                className="hover:cursor-pointer"
                                onClick={() => deleteItem([provider.name])}
                            >
                                <Button onPress={()=>{setEdit(true); setEditValue({name: provider.name, path: provider.path, remote: provider.remote}); onOpen()}}>
                                    <span className={`text-md w-28 text-wrap font-black select-none`}>
                                        {provider.name}
                                    </span>
                                </Button>
                            </Badge>
                            <Checkbox value={provider.name} />
                        </div>
                    ))}
                </div>
            </CheckboxGroup>
        </ScrollShadow>
    )
}