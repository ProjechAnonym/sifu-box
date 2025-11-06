import { useState, useMemo } from "react";
import { ScrollShadow } from "@heroui/scroll-shadow";
import { Badge } from "@heroui/badge";
import { Button } from "@heroui/button";
import { Checkbox, CheckboxGroup } from "@heroui/checkbox";
import { Input } from "@heroui/input";
import { Provider } from "@/types/setting/provider";
export default function ProviderLayout(props: {providers: Array<Provider>}) {
    const { providers } = props;
    const [filtered_providers, setFilteredProviders] = useState(providers);
    const [checked, setChecked] = useState<Array<string>>([]);
    const [search, setSearch] = useState("");
    useMemo(() => {
        providers &&
        search &&
        setFilteredProviders(
            providers.filter((provider) => provider.name.includes(search))
        );
        providers && search === "" && setFilteredProviders(providers);
    }, [search, providers]);
    return (
        <ScrollShadow className="w-full h-1/2">
            <header className="flex flex-wrap gap-2 p-1 items-center">
                <Input
                    size="sm"
                    label={<span className="text-md font-black">机场</span>}
                    variant="underlined"
                    value={search}
                    onValueChange={setSearch}
                    className="w-24"
                />
                <Button color="danger" size="sm" variant="shadow" onPress={()=>console.log(checked)}>
                    <span className="text-xl font-black">删除</span>
                </Button>
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
                                onClick={() => console.log(provider.name)}
                            >
                                <Button onPress={()=>console.log("5")}>
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