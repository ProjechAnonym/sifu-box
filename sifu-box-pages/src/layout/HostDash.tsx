import SetHost from "@/components/setting/setHost";
import { HostValue } from "@/types/host";
export default function HostDash(props: {
  secret: string;
  dark: boolean;
  templates: Array<{ Name: string; Template: Object }> | null;
  hosts: Array<HostValue> | null;
  setUpdateHosts: (Updatehosts: boolean) => void;
}) {
  const { secret, dark, hosts, templates, setUpdateHosts } = props;

  return (
    <header className="p-2 flex flex-wrap gap-2 items-center">
      <SetHost
        secret={secret}
        hosts={hosts}
        setUpdateHosts={setUpdateHosts}
        dark={dark}
        templates={templates}
      />
    </header>
  );
}
