package sqlutil

import (
	"strings"
)

func AddTenantCheck(tenantID int, filterExp string) string {
	if tenantID > 0 && len(filterExp) > 0 {
		return "tenant_id = ? AND " + filterExp
	}
	if tenantID > 0 {
		return "tenant_id = ?"
	}
	return ""
}

func FmtSQL(selectExp, whereExp, sortExp string, limit, offset int) string {
	// pre-over-allocate cause we can ^.^
	out := strings.Builder{}
	out.Grow(len(" WHERE "+whereExp) + len(" ORDER BY "+sortExp) + 20) // 20 == len (" limit ? + offset ? ")

	out.WriteString(selectExp)
	if len(whereExp) > 0 {
		out.WriteString(" WHERE " + whereExp)
	}
	if len(sortExp) > 0 {
		out.WriteString(" ORDER BY " + sortExp)
	}
	if limit > 0 {
		out.WriteString(" LIMIT ?")
	}
	if offset > 0 {
		out.WriteString(" OFFSET ?")
	}
	return out.String()
}

func FmtSQLArgs(tenantID, limit, offset int, args []interface{}) []interface{} {
	// insert at tenantID front because order matters and it will be first
	if tenantID > 0 {
		args = append([]interface{}{tenantID}, args...)
	}
	if limit > 0 {
		args = append(args, limit)
	}
	if offset > 0 {
		args = append(args, offset)
	}
	return args
}

func FmtListQuery(tenantID int, listSQL, filterExp, sortExp string, limit, offset int) string {
	whereExp := AddTenantCheck(tenantID, filterExp)
	sql := FmtSQL(listSQL, whereExp, sortExp, limit, offset)
	return sql
}
