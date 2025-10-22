
import { useState, useRef, useMemo } from "react";
import { Button } from "@heroui/button";

export default function NavBar(props: { groups: Array<string> }) {
  const { groups } = props;
  const block_div = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(true);
  const [nav_item, setNavItem] = useState<string>("");
  useMemo(() => {
    const element = document.getElementById(nav_item);
    element && element.scrollIntoView({ behavior: "smooth" });
  }, [nav_item]);
  return (
    <div className="absolute top-1/2 left-4 z-50 w-32">
      <Button
        onPress={() => {
          setIsOpen(!isOpen);
          !isOpen &&
            block_div.current &&
            (block_div.current.style.display = "block");
        }}
        color="primary"
        size="sm"
        isIconOnly
      >
        <i className={`bi bi-${isOpen ? "x" : "list"} text-2xl`} />
      </Button>
      <div
        ref={block_div}
        onAnimationEnd={() =>
          isOpen
            ? block_div.current && (block_div.current.style.display = "block")
            : block_div.current && (block_div.current.style.display = "none")
        }
        className={`bottom-0 right-0 w-full h-full bg-gradient-to-tr from-[#6c89fc] to-[#ff9f9f] p-2 rounded-lg ${
          isOpen
            ? "translate-y-0 animate-expand-open"
            : "-translate-y-1/2 animate-expand-close"
        }`}
      >
        {groups.map((group) => (
          <Button
            key={group}
            size="sm"
            onPress={() => setNavItem(group)}
            className="my-0.5"
            variant="light"
          >
            <span className="text-lg font-black">{group}</span>
          </Button>
        ))}
      </div>
    </div>
  );
}
