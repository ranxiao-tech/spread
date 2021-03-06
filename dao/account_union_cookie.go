package dao

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/spread/global"
	"time"
	"fmt"
)

func GetUnionCookies(platformId, accId int64, platformAcc string, timeOut time.Duration) (string, string, string, error) {
	m, err := GetUnionCookie(platformId, accId, platformAcc, timeOut)
	if err != nil {
		return "", "", "", err
	}
	return convert.ToString(m["cookies"]),convert.ToString(m["local_storage"]),convert.ToString(m["sessions"]), nil
}

func SetUnionCookies(domain, cookies, localStorage, sessions string, platformId, accId int64, platformAcc string) (bool, error) {
	if global.Acc == nil {
		return false, nil
	}
	m, err := GetUnionCookie(platformId, accId, platformAcc, -1)
	if err != nil {
		return false, err
	}
	var x int64
	var err2 error
	if m == nil {
		//修改原有Cookie
		sql := `INSERT INTO account_union_cookie(id,account_id,domain,cookies,ut_time,ct_time,status,platform_account,platform_id,local_storage,sessions) values(?,?,?,?,?,?,1,?,?,?)`
		nt := time.Now().Format("2006-01-02 15:04:05")
		x, err2 = global.SLDB.Insert(sql, utils.ID(), accId, domain, cookies, nt, nt, platformAcc, platformId, localStorage, sessions)
	} else {
		sql := `UPDATE account_union_cookie SET  cookies=?,local_storage=?,sessions=?,ut_time=? WHERE platform_id=? AND account_id=? AND
		platform_account=?`
		nt := time.Now().Format("2006-01-02 15:04:05")
		x, err2 = global.SLDB.Update(sql, cookies, localStorage, sessions, nt, platformId, accId, platformAcc)
	}
	if err2 != nil {
		return false, err
	}
	return x > 0, err
}

func GetUnionCookie(platformId, accId int64, platformAcc string, timeOut time.Duration) (map[string]interface{}, error) {
	if global.Acc == nil {
		return nil, nil
	}
	var timeWhere string

	if timeOut == -2 {
		return nil, nil
	}
	if timeOut != -1 && timeOut != -3 {
		nt := utils.FormatTime(time.Now().Add(-(time.Minute * timeOut)))
		timeWhere = fmt.Sprintf(" AND ut_time>'%s'", nt)
	}
	sql := fmt.Sprintf(`SELECT cookies,local_storage,sessions FROM account_union_cookie WHERE platform_id=? AND account_id=? AND platform_account=? %s ORDER BY	ut_time DESC LIMIT 1`, timeWhere)
	m, err := global.SLDB.Query(sql, platformId, accId, platformAcc)
	if err != nil {
		return nil, err
	}
	if len(m) <= 0 {
		return nil, nil
	}
	return m[0], nil
}
