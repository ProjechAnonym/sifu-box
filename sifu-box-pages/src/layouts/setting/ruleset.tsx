import { useState, useMemo, useRef } from "react";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Badge } from "@heroui/badge";
import { Button } from "@heroui/button";
import { Checkbox, CheckboxGroup } from "@heroui/checkbox";
import { Input } from "@heroui/input";
import { useDisclosure } from "@heroui/modal";
import toast from "react-hot-toast";
import AddRuleset from "@/components/modal/ruleset";
import { DeleteRuleset, AddRulesetFiles } from "@/utils/configuration/ruleset";
import { RuleSet, DEFAULT_RULESET } from "@/types/setting/ruleset";
export default function RulesetLayout(props: {rulesets: Array<RuleSet>, theme: string, token: string, setUpdate: (update: boolean) => void}) {
    const { rulesets, theme, token, setUpdate } = props;
    const { isOpen, onOpen, onClose } = useDisclosure();
    const file_input = useRef<HTMLInputElement>(null);
    const [filtered_rulesets, setFilteredRulesets] = useState(rulesets);
    const [checked, setChecked] = useState<Array<string>>([]);
    const [search, setSearch] = useState("");
    const [edit, setEdit] = useState(false);
    const [edit_value, setEditValue] = useState<{name: string, path: string, remote: boolean, binary: boolean, download_detour: string, update_interval: string}>(DEFAULT_RULESET);
    const [files, setFiles] = useState<FileList | null>(null);
    useMemo(() => {
        rulesets &&
        search &&
        setFilteredRulesets(
            rulesets.filter((ruleset) => ruleset.name.includes(search))
        );
        rulesets && search === "" && setFilteredRulesets(rulesets);
    }, [search, rulesets]);
    const deleteItem = (value: Array<string>) => toast.promise(DeleteRuleset(token, value), {
        loading: "正在删除所选规则集...",
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
    const addFiles = (files: FileList) => toast.promise(AddRulesetFiles(token, files), {
        loading: "正在添加所选规则集...",
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
    const openModal = (ruleset: {name: string, path: string, remote: boolean, binary: boolean, download_detour: string, update_interval: string}, edit: boolean) => {
        setEdit(edit); 
        setEditValue(ruleset);
        onOpen()
    }
    return (
        <ScrollShadow className="w-full h-1/2">
            <AddRuleset edit={edit} isOpen={isOpen} onClose={onClose} theme={theme} token={token} setUpdate={setUpdate} initial_value={edit_value} />
            <header className="flex flex-wrap gap-1 p-1 items-center">
                <Input
                    size="sm"
                    label={<span className="text-md font-black">规则集</span>}
                    variant="underlined"
                    value={search}
                    onValueChange={setSearch}
                    className="w-24"
                />
                <Button color="danger" size="sm" variant="shadow" onPress={() => deleteItem(checked)}>
                    <span className="text-xl font-black">删除</span>
                </Button>
                <Button color="primary" size="sm" variant="shadow" onPress={() => openModal(DEFAULT_RULESET, false)}>
                    <span className="text-xl font-black">添加</span>
                </Button>
                <form className="flex flex-row gap-1" 
                    onSubmit={e => {
                        e.preventDefault()
                        files ? addFiles(files) : toast.error("请选择文件")
                    }}
                >
                    <Badge content={files ? files.length : 0} color="success" placement="bottom-right">
                        <label htmlFor="file-upload-ruleset" className="px-2 bg-primary h-8 rounded-md hover:cursor-pointer hover:bg-opacity-85 transition-all text-white">
                            <i className="bi bi-filetype-raw text-2xl" />
                        </label>
                    </Badge>
                    <input
                        type="file"
                        id="file-upload-ruleset"
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
                    {filtered_rulesets && filtered_rulesets.map((ruleset) => (
                        <div key={ruleset.name} className="flex flex-row justify-center h-fit gap-1">
                            <Badge
                                content={<i className="bi bi-trash-fill" />}
                                placement="bottom-left"
                                color="danger"
                                shape="rectangle"
                                className="hover:cursor-pointer"
                                onClick={() => deleteItem([ruleset.name])}
                            >
                                <Button onPress={()=>openModal(ruleset, true)} isDisabled={!ruleset.remote}>
                                    <span className={`text-md w-28 text-wrap font-black select-none`}>
                                        {ruleset.name}
                                    </span>
                                </Button>
                            </Badge>
                            <Checkbox value={ruleset.name} />
                        </div>
                    ))}
                </div>
            </CheckboxGroup>
        </ScrollShadow>
    )
}