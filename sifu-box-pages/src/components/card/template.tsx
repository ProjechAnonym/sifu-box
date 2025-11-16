import { Card, CardBody, CardFooter, CardHeader } from "@heroui/card"
import { Popover, PopoverContent, PopoverTrigger } from "@heroui/popover"
import { Button } from "@heroui/button"
import { Checkbox } from "@heroui/checkbox"
import { ScrollShadow } from "@heroui/scroll-shadow"
export default function TemplateCard(props: {template: {name: string, [key: string]: any}; theme: string; setDelete: (name: Array<string>) => void; token: string; setTemplate: (template: string) => void; editTemplate: (template: {name: string, [key: string]: any}) => void;}){ 
    const {template, theme, setDelete, setTemplate, editTemplate} = props
    return (
        <Card shadow="sm" key={template.name}>
            <CardHeader className="justify-between">
                <Popover shadow="sm" classNames={{content: `${theme} bg-content1 text-foreground`}}>
                    <PopoverTrigger>
                        <Button size="sm" variant="shadow">
                        <span className={`text-xl font-black`}>{template.name}</span>
                        </Button>
                    </PopoverTrigger>
                    <PopoverContent>
                        <p className="text-md font-black w-36 p-1">
                            是否将"{template.name}"模板设置为活动模板
                        </p>
                        <p className="w-full justify-end flex p-1">
                        <Button
                            size="sm"
                            color="primary"
                            variant="shadow"
                            onPress={()=>setTemplate(template.name)}
                        >
                            <span className="text-xl font-black">确认</span>
                        </Button>
                        </p>
                    </PopoverContent>
                </Popover>
                <Checkbox value={template.name} />
            </CardHeader>
            <CardBody>
                <ScrollShadow key={template.name} className="w-full h-40">{JSON.stringify(template, null, 4)}</ScrollShadow>
            </CardBody>
            <CardFooter className="justify-end gap-2">
                <Button color="danger" size="sm" variant="shadow" onPress={()=>setDelete([template.name])}>
                    <span className="text-xl font-black">删除</span>
                </Button>
                <Button color="primary" size="sm" variant="shadow" onPress={() => editTemplate(template)}>
                <span className="text-xl font-black">修改</span>
                </Button>
            </CardFooter>
        </Card>
    )
}