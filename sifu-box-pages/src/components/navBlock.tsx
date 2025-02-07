import { useState, useRef, useMemo } from "react";
import { Button } from "@heroui/button";

export default function NavBlock(props: { groups: Array<string> }) {
  const { groups } = props;
  const blockDiv = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [navItem, setNavItem] = useState<string>("");
  useMemo(() => {
    const element = document.getElementById(navItem);
    element && element.scrollIntoView({ behavior: "smooth" });
  }, [navItem]);
  return (
    <div className="absolute top-1/2 left-4 z-50 w-32">
      <Button
        onPress={() => {
          setIsOpen(!isOpen);
          !isOpen &&
            blockDiv.current &&
            (blockDiv.current.style.display = "block");
        }}
        color="primary"
        size="sm"
        isIconOnly
      >
        <i className={`bi bi-${isOpen ? "x" : "list"} text-2xl`} />
      </Button>
      <div
        ref={blockDiv}
        onAnimationEnd={() =>
          isOpen
            ? blockDiv.current && (blockDiv.current.style.display = "block")
            : blockDiv.current && (blockDiv.current.style.display = "none")
        }
        className={`bottom-0 right-0 w-full h-full bg-gradient-to-tr from-[#6c89fc] to-[#ff9f9f] p-2 rounded-lg ${
          isOpen
            ? "translate-y-0 animate-showIn_normal"
            : "-translate-y-1/2 animate-showOut_normal"
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
