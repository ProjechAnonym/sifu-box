import CustomTable from "../table";
import { initRuleLabels, ruleColumns } from "@/types/singbox/rules";

export default function Rules(props: {
  theme: string;
  rules: Array<ruleColumns>;
}) {
  const { rules, theme } = props;

  return (
    <div>
      <CustomTable
        theme={theme}
        data={{ labels: initRuleLabels, values: rules }}
        defaultSearchField="payload"
      />
    </div>
  );
}
