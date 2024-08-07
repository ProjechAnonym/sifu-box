import { useRef, useEffect } from "react";
import { Tooltip, Link, Image } from "@nextui-org/react";
import "@/assets/ali-iconfont-social/iconfont.css";
import alipay from "@/assets/pictures/alipay.jpg";
import wechat from "@/assets/pictures/wechat.jpg";
export default function FooterBar(props: {
  height_callback: (height: number) => void;
}) {
  const ref = useRef<HTMLDivElement>(null);
  const { height_callback } = props;
  useEffect(() => {
    ref.current && height_callback(ref.current.clientHeight);
  }, [ref]);
  return (
    <footer className="w-full h-15 text-center space-x-2 py-2" ref={ref}>
      <p>
        <Tooltip content="sifulin's blog" placement="top" offset={5}>
          <Link
            className="text-foreground"
            href="https://vercel-blog.sifulin.top"
            isExternal
            isBlock
          >
            <i className="iconfont icon-boke1 text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip content="Github" placement="top" offset={5}>
          <Link
            className="text-foreground"
            href="https://github.com/ProjechAnonym"
            isExternal
            isBlock
          >
            <i className="iconfont icon-github text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip content="Youtube" placement="top" offset={5}>
          <Link
            className="text-foreground"
            href="https://www.youtube.com/channel/UCXiiRClqjDLrqzMbq2Kqb4A"
            isExternal
            isBlock
          >
            <i className="iconfont icon-youtube text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip content="Telegram" placement="top" offset={5}>
          <Link
            className="text-foreground"
            href="https://t.me/+5yh2rgXjWBlmMDk1"
            isExternal
            isBlock
          >
            <i className="iconfont icon-telegram-original text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip content="QQ" placement="top" offset={5}>
          <Link
            className="text-foreground"
            href="https://qm.qq.com/cgi-bin/qm/qr?authKey=GAwy5K83J0XtzMdoajfHMWauVzqwzUawF%2F8vVlKSoGncd9InsiqRsssT1ybQH1tY&k=X9BMdHD7h8Qk1FgX9T2aJXhoFOZ6j_0n&noverify=0"
            isExternal
            isBlock
          >
            <i className="iconfont icon-qq text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip content="bilibili" placement="top" offset={5}>
          <Link
            className="text-foreground"
            href="https://space.bilibili.com/8337954"
            isExternal
            isBlock
          >
            <i className="iconfont icon-bilibili text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip
          content={
            <div className="flex flex-row gap-1">
              <Image width={64} height={64} alt="alipay" src={alipay} />
              <Image width={64} height={64} alt="wechat" src={wechat} />
            </div>
          }
          placement="top"
          offset={5}
        >
          <Link className="text-foreground" isBlock href="#">
            <i className="iconfont icon-dashang1 text-2xl" />
          </Link>
        </Tooltip>
      </p>
      <p className="font-semibold">developed by sifulin</p>
    </footer>
  );
}
