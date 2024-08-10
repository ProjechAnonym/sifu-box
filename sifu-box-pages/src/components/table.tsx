import { useMemo, useState } from "react";
import {
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  getKeyValue,
  SortDescriptor,
  DropdownTrigger,
  Dropdown,
  DropdownMenu,
  DropdownItem,
  Button,
  Selection,
  Input,
  Pagination,
  Select,
  SelectItem,
} from "@nextui-org/react";

export default function MyTable(props: {
  data: {
    labels: Array<{
      key: string;
      label: string;
      allowSort: boolean;
      initShow: boolean;
    }>;
    values: Array<{ key: string; [addtionalProp: string]: any }>;
  };
  dark?: boolean;
  defaultSearchField?: string;
  rowsPerPage?: number;
  tableHeight?: string;
  dropdownHeight?: string;
}) {
  const {
    data,
    dark = true,
    defaultSearchField = "",
    rowsPerPage = 10,
    tableHeight = "h-96",
    dropdownHeight = "h-48",
  } = props;
  const { labels, values } = data;
  const [filterValue, setFilterValue] = useState("");
  const [searchField, setSearchField] = useState<Selection>(
    new Set([defaultSearchField])
  );
  const [sortDescriptor, setSortDescriptor] = useState<SortDescriptor>({
    column: labels[0].key,
    direction: "ascending",
  });
  const [visibleLabels, setVisibleLabels] = useState<Selection>(
    new Set(labels.filter((label) => label.initShow).map((label) => label.key))
  );
  const [page, setPage] = useState(1);
  const showLabels = useMemo(() => {
    return labels.filter((label) =>
      Array.from(visibleLabels).includes(label.key)
    );
  }, [visibleLabels]);
  const changeSearchField = useMemo(
    () => Array.from(searchField)[0],
    [searchField]
  );
  const filteredValues = useMemo(
    () =>
      filterValue !== "" && searchField
        ? values.filter(
            (value) =>
              value[changeSearchField] &&
              (value[changeSearchField] as string)
                .toLowerCase()
                .includes(filterValue.toLowerCase())
          )
        : values,
    [values, filterValue]
  );
  const pages = Math.ceil(filteredValues.length / rowsPerPage);
  const showValues = useMemo(() => {
    const start = (page - 1) * rowsPerPage;
    const end = start + rowsPerPage;
    return filteredValues.slice(start, end);
  }, [page, filteredValues, rowsPerPage]);
  const sortValues = useMemo(
    () =>
      showValues.sort((a, b) => {
        const first = a[sortDescriptor.column!];
        const second = b[sortDescriptor.column!];
        const cmp =
          (Number(first) || first) < (Number(second) || second) ? -1 : 1;
        return sortDescriptor.direction === "descending" ? -cmp : cmp;
      }),
    [sortDescriptor, showValues]
  );
  return (
    <>
      <Table
        topContentPlacement="outside"
        bottomContentPlacement="outside"
        classNames={{ wrapper: tableHeight }}
        topContent={
          <header className="flex flex-row items-center gap-x-2 justify-center">
            {searchField && (
              <Input
                size="sm"
                placeholder="search ..."
                startContent={<i className="bi bi-search" />}
                value={filterValue}
                onValueChange={setFilterValue}
                isClearable
                onClear={() => setFilterValue("")}
              />
            )}
            <Select
              size="sm"
              classNames={{
                popoverContent: `${
                  dark ? "sifudark bg-zinc-800" : "sifulight bg-slate-100"
                } text-foreground`,
              }}
              aria-label="Search Columns"
              selectedKeys={searchField}
              onSelectionChange={setSearchField}
              placeholder="Search Columns"
            >
              {labels.map((label) => (
                <SelectItem key={`${label.key}`} textValue={label.label}>
                  <span className="font-black">{label.label}</span>
                </SelectItem>
              ))}
            </Select>
            <Dropdown
              className={`${
                dark ? "sifudark bg-zinc-800" : "sifulight bg-slate-100"
              } text-foreground`}
            >
              <DropdownTrigger className="hidden sm:flex">
                <Button
                  endContent={<i className="bi bi-chevron-down" />}
                  variant="flat"
                  size="sm"
                >
                  <span className="font-black">列名</span>
                </Button>
              </DropdownTrigger>
              <DropdownMenu
                disallowEmptySelection
                aria-label="Table Columns"
                closeOnSelect={false}
                selectedKeys={visibleLabels}
                selectionMode="multiple"
                onSelectionChange={setVisibleLabels}
                classNames={{ base: `${dropdownHeight} overflow-y-auto` }}
              >
                {labels.map((label) => (
                  <DropdownItem
                    key={label.key}
                    className="capitalize"
                    textValue={label.label}
                  >
                    <span className="font-black">{label.label}</span>
                  </DropdownItem>
                ))}
              </DropdownMenu>
            </Dropdown>
          </header>
        }
        bottomContent={
          pages > 1 && (
            <footer className="flex justify-center">
              <Pagination
                siblings={2}
                isCompact
                total={pages}
                showControls
                loop
                page={page}
                onChange={setPage}
              />
            </footer>
          )
        }
        aria-label="Example table with dynamic content"
        sortDescriptor={sortDescriptor}
        onSortChange={setSortDescriptor}
      >
        <TableHeader columns={showLabels}>
          {(column) => (
            <TableColumn key={column.key} allowsSorting={column.allowSort}>
              {column.label}
            </TableColumn>
          )}
        </TableHeader>
        <TableBody items={sortValues} emptyContent={"No rows to display."}>
          {sortValues.length === 0
            ? []
            : (item) => (
                <TableRow key={item.key}>
                  {(columnKey) => (
                    <TableCell>{getKeyValue(item, columnKey)}</TableCell>
                  )}
                </TableRow>
              )}
        </TableBody>
      </Table>
    </>
  );
}
