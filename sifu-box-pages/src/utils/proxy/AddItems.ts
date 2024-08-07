import axios from "axios";
import {
  ProviderData,
  ProviderValue,
  RulesetLocalData,
  RulesetRemoteData,
  RulesetValue,
} from "@/types/proxy";
export async function AddFile(secret: string, files: FileList) {
  const formData = new FormData();
  for (let i = 0; i < files.length; i++) {
    files.item(i) &&
      formData.append("files", files.item(i)!, files.item(i)!.name);
  }
  try {
    const res = await axios.post("/api/proxy/files", formData, {
      headers: { Authorization: secret },
    });
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
export async function AddProxy(
  secret: string,
  proxy: {
    providers: Array<ProviderValue>;
    rulesets: Array<RulesetValue>;
  }
) {
  const newProviders = proxy.providers.map((provider): ProviderData => {
    const { id, ...providerData } = provider;
    return providerData;
  });
  const newRulesets = proxy.rulesets.map(
    (ruleset): RulesetRemoteData | RulesetLocalData => {
      if (ruleset.type === "local") {
        const { id, update_interval, download_detour, url, ...rulesetData } =
          ruleset;
        return rulesetData;
      } else {
        const { id, path, ...rulesetData } = ruleset;
        return rulesetData;
      }
    }
  );
  try {
    const res = await axios.post(
      "/api/proxy/add",
      { providers: newProviders, rulesets: newRulesets },
      {
        headers: { Authorization: secret },
      }
    );
    return res.status === 200;
  } catch (e) {
    console.error(e);
    throw e;
  }
}
