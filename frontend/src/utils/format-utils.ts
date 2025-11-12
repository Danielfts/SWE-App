export const formatColombianDateTime = (isoString: string): string => {
  const date = new Date(isoString);
  return date.toLocaleString('es-CO', {
    timeZone: 'America/Bogota',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  });
};

export function formatAsMoney(value: number) {
  let newStr = "";
  const [intPart, decPart] = value.toString().split(".")
  const oldStr = intPart!.toString().split("").reverse();
  let counter = 0;
  for (let i = 0; i < oldStr.length; i++) {
    const ch = oldStr[i]
    newStr = ch + newStr
    counter++;
    if (counter === 3 && i < oldStr.length - 1) {
      newStr = ',' + newStr;
      counter = 0;
    }
  }
  newStr = "$" + newStr + '.' + decPart;
  return newStr;
}

export function getDelta(from: number, to: number): string {
  const percentage = ((to - from) / from) * 100;
  return percentage.toFixed(1) + '%';
}

export function compareDecimals(from: any, to: any): number {
  const fromN = parseFloat(from);
  const toN = parseFloat(to);
  if (toN > fromN) return 1
  else if (toN < fromN) return -1
  else return 0
}

export const getSortChar = (label: string, sortby: string, sortorder: boolean) => {
  if (label === sortby) {
    if (sortorder) {
      return "^"
    }
    return "v"
  }
  return "-"
}