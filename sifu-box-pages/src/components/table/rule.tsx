import MyTable from "@/layouts/mytable";
import { init_rule_labels, ruleColumns } from "@/types/singbox/rules";

export default function Rules(props: {
  theme: string;
  rules: Array<ruleColumns>;
}) {
  const { rules, theme } = props;

  return (
    <div>
      <MyTable
        theme={theme}
        data={{ labels: init_rule_labels, values: rules }}
        defaultSearchField="payload"
      />
    </div>
  );
}