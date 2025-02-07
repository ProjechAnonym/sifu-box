import { Tooltip } from "@heroui/tooltip";
import { Link } from "@heroui/link";
import { Image } from "@heroui/image";
import "@/assets/ali-iconfont-social/iconfont.css";
import alipay from "@/assets/pictures/alipay.jpg";
import wechat from "@/assets/pictures/wechat.jpg";
export default function Footer(props: { theme: string }) {
  const { theme } = props;
  return (
    <footer className="flex flex-col items-center justify-center w-full h-16 p-2">
      <p>
        <Tooltip
          content={
            <span className={`${theme} text-foreground`}>sifulin's blog</span>
          }
          placement="top"
          offset={5}
          classNames={{
            content: [`${theme} bg-content1`],
          }}
        >
          <Link
            className="text-foreground"
            href="https://vercel-blog.sifulin.top"
            isExternal
            isBlock
          >
            <i className="iconfont icon-boke1 text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip
          content={<span className={`${theme} text-foreground`}>Github</span>}
          placement="top"
          offset={5}
          classNames={{
            content: [`${theme} bg-content1`],
          }}
        >
          <Link
            className="text-foreground"
            href="https://github.com/ProjechAnonym"
            isExternal
            isBlock
          >
            <i className="iconfont icon-github text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip
          content={<span className={`${theme} text-foreground`}>Youtube</span>}
          placement="top"
          offset={5}
          classNames={{
            content: [`${theme} bg-content1`],
          }}
        >
          <Link
            className="text-foreground"
            href="https://www.youtube.com/channel/UCXiiRClqjDLrqzMbq2Kqb4A"
            isExternal
            isBlock
          >
            <i className="iconfont icon-youtube text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip
          content={<span className={`${theme} text-foreground`}>Telegram</span>}
          placement="top"
          offset={5}
          classNames={{
            content: [`${theme} bg-content1`],
          }}
        >
          <Link
            className="text-foreground"
            href="https://t.me/+5yh2rgXjWBlmMDk1"
            isExternal
            isBlock
          >
            <i className="iconfont icon-telegram-original text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip
          content={<span className={`${theme} text-foreground`}>QQ</span>}
          placement="top"
          offset={5}
          classNames={{
            content: [`${theme} bg-content1`],
          }}
        >
          <Link
            className="text-foreground"
            href="https://qm.qq.com/cgi-bin/qm/qr?authKey=GAwy5K83J0XtzMdoajfHMWauVzqwzUawF%2F8vVlKSoGncd9InsiqRsssT1ybQH1tY&k=X9BMdHD7h8Qk1FgX9T2aJXhoFOZ6j_0n&noverify=0"
            isExternal
            isBlock
          >
            <i className="iconfont icon-qq text-2xl" />
          </Link>
        </Tooltip>
        <Tooltip
          content={<span className={`${theme} text-foreground`}>bilibili</span>}
          placement="top"
          offset={5}
          classNames={{
            content: [`${theme} bg-content1`],
          }}
        >
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
            <div className="flex flex-row gap-1 py-1">
              <Image width={64} height={64} alt="alipay" src={alipay} />
              <Image width={64} height={64} alt="wechat" src={wechat} />
            </div>
          }
          placement="top"
          offset={5}
          classNames={{
            content: [`${theme} bg-content1`],
          }}
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
