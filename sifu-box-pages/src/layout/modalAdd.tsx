import { useState, useEffect } from "react";
import { ScrollShadow, Button } from "@nextui-org/react";
import AddUrl from "@/components/proxy/addProvider";
import AddRuleset from "@/components/proxy/addRuleset";
import { cloneDeep } from "lodash";
import { ProviderValue, RulesetValue } from "@/types/proxy";

export default function ModalAdd(props: {
  onProxySubmit: (proxy: {
    providers: Array<ProviderValue>;
    rulesets: Array<RulesetValue>;
  }) => void;
}) {
  const { onProxySubmit } = props;
  const [providers, setProviders] = useState<Array<ProviderValue>>([]);
  const [rulesets, setRulesets] = useState<Array<RulesetValue>>([]);
  useEffect(() => {
    onProxySubmit({ providers: providers, rulesets: rulesets });
  }, [providers, rulesets]);
  return (
    <div className="flex flex-col gap-y-2">
      <header className={`text-foreground flex flex-row gap-x-2 items-center`}>
        机场链接
        <Button
          isIconOnly
          size="sm"
          onPress={() => {
            const newProviders = cloneDeep(providers);
            newProviders.push({
              id: 0,
              path: "",
              proxy: false,
              name: "",
              remote: true,
            });
            setProviders(newProviders);
          }}
        >
          <span>
            <i className="bi bi-plus" />
          </span>
        </Button>
        <Button
          isIconOnly
          size="sm"
          onPress={() => {
            const newProviders = cloneDeep(providers);
            newProviders.pop();
            setProviders(newProviders);
          }}
        >
          <span>
            <i className="bi bi-dash" />
          </span>
        </Button>
      </header>
      <ScrollShadow className="h-36 flex flex-col gap-y-1">
        {providers.map((_, i) => (
          <AddUrl
            key={`providers-${i}`}
            onSubmit={(provider) => {
              const newProviders = cloneDeep(providers);
              newProviders[i] = provider;
              setProviders(newProviders);
            }}
          />
        ))}
      </ScrollShadow>
      <header className={`text-foreground flex flex-row gap-x-2 items-center`}>
        规则集链接
        <Button
          isIconOnly
          size="sm"
          onPress={() => {
            const new_rulesets = cloneDeep(rulesets);
            new_rulesets.push({
              tag: "",
              id: 0,
              url: "",
              dnsRule: "",
              path: "",
              format: "",
              type: "",
              china: false,
              download_detour: "",
              update_interval: "1d",
              label: "",
            });
            setRulesets(new_rulesets);
          }}
        >
          <span>
            <i className="bi bi-plus" />
          </span>
        </Button>
        <Button
          isIconOnly
          size="sm"
          onPress={() => {
            const newRulesets = cloneDeep(rulesets);
            newRulesets.pop();
            setRulesets(newRulesets);
            onProxySubmit({ providers: providers, rulesets: newRulesets });
          }}
        >
          <span>
            <i className="bi bi-dash" />
          </span>
        </Button>
      </header>
      <ScrollShadow className="h-36 flex flex-col gap-y-2">
        {rulesets.map((_, i) => (
          <AddRuleset
            key={`rulesets-${i}`}
            onSubmit={(ruleset) => {
              const newRulesets = cloneDeep(rulesets);
              newRulesets[i] = ruleset;
              setRulesets(newRulesets);
            }}
          />
        ))}
      </ScrollShadow>
    </div>
  );
}
