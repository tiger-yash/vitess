/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package evalengine

import (
	"math"
	"time"

	"vitess.io/vitess/go/hack"
	"vitess.io/vitess/go/mysql/collations"
	"vitess.io/vitess/go/mysql/datetime"
	"vitess.io/vitess/go/mysql/decimal"
	"vitess.io/vitess/go/sqltypes"
	querypb "vitess.io/vitess/go/vt/proto/query"
)

var SystemTime = time.Now

type (
	builtinNow struct {
		CallExpr
		utc      bool
		onlyTime bool
		prec     uint8
	}

	builtinSysdate struct {
		CallExpr
		prec uint8
	}

	builtinCurdate struct {
		CallExpr
	}

	builtinUtcDate struct {
		CallExpr
	}

	builtinDateFormat struct {
		CallExpr
		collate collations.ID
	}

	builtinDate struct {
		CallExpr
	}

	builtinDayOfMonth struct {
		CallExpr
	}

	builtinDayOfWeek struct {
		CallExpr
	}

	builtinDayOfYear struct {
		CallExpr
	}

	builtinFromUnixtime struct {
		CallExpr
		collate collations.ID
	}

	builtinHour struct {
		CallExpr
	}

	builtinMakedate struct {
		CallExpr
	}

	builtinMaketime struct {
		CallExpr
	}

	builtinMicrosecond struct {
		CallExpr
	}

	builtinMinute struct {
		CallExpr
	}

	builtinMonth struct {
		CallExpr
	}

	builtinMonthName struct {
		CallExpr
		collate collations.ID
	}

	builtinQuarter struct {
		CallExpr
	}

	builtinSecond struct {
		CallExpr
	}

	builtinTime struct {
		CallExpr
	}

	builtinUnixTimestamp struct {
		CallExpr
	}

	builtinWeek struct {
		CallExpr
	}

	builtinWeekDay struct {
		CallExpr
	}

	builtinWeekOfYear struct {
		CallExpr
	}

	builtinYear struct {
		CallExpr
	}

	builtinYearWeek struct {
		CallExpr
	}

	builtinDateMath struct {
		CallExpr
		sub     bool
		unit    datetime.IntervalType
		collate collations.ID
	}
)

var _ Expr = (*builtinNow)(nil)
var _ Expr = (*builtinSysdate)(nil)
var _ Expr = (*builtinCurdate)(nil)
var _ Expr = (*builtinUtcDate)(nil)
var _ Expr = (*builtinDateFormat)(nil)
var _ Expr = (*builtinDate)(nil)
var _ Expr = (*builtinDayOfMonth)(nil)
var _ Expr = (*builtinDayOfWeek)(nil)
var _ Expr = (*builtinDayOfYear)(nil)
var _ Expr = (*builtinHour)(nil)
var _ Expr = (*builtinFromUnixtime)(nil)
var _ Expr = (*builtinMakedate)(nil)
var _ Expr = (*builtinMaketime)(nil)
var _ Expr = (*builtinMicrosecond)(nil)
var _ Expr = (*builtinMinute)(nil)
var _ Expr = (*builtinMonth)(nil)
var _ Expr = (*builtinMonthName)(nil)
var _ Expr = (*builtinQuarter)(nil)
var _ Expr = (*builtinSecond)(nil)
var _ Expr = (*builtinTime)(nil)
var _ Expr = (*builtinUnixTimestamp)(nil)
var _ Expr = (*builtinWeek)(nil)
var _ Expr = (*builtinWeekDay)(nil)
var _ Expr = (*builtinWeekOfYear)(nil)
var _ Expr = (*builtinYear)(nil)
var _ Expr = (*builtinYearWeek)(nil)

func (call *builtinNow) eval(env *ExpressionEnv) (eval, error) {
	now := env.time(call.utc)
	if call.onlyTime {
		buf := datetime.Time_hh_mm_ss.Format(now, call.prec)
		return newEvalRaw(sqltypes.Time, buf, collationBinary), nil
	} else {
		buf := datetime.DateTime_YYYY_MM_DD_hh_mm_ss.Format(now, call.prec)
		return newEvalRaw(sqltypes.Datetime, buf, collationBinary), nil
	}
}

func (call *builtinNow) typeof(_ *ExpressionEnv, _ []*querypb.Field) (sqltypes.Type, typeFlag) {
	if call.onlyTime {
		return sqltypes.Time, 0
	}
	return sqltypes.Datetime, 0
}

func (call *builtinNow) compile(c *compiler) (ctype, error) {
	var format *datetime.Strftime
	var t sqltypes.Type

	if call.onlyTime {
		format = datetime.Time_hh_mm_ss
		t = sqltypes.Time
	} else {
		format = datetime.DateTime_YYYY_MM_DD_hh_mm_ss
		t = sqltypes.Datetime
	}
	c.asm.Fn_Now(t, format, call.prec, call.utc)
	return ctype{Type: t, Col: collationBinary}, nil
}

func (call *builtinNow) constant() bool {
	return false
}

func (call *builtinSysdate) eval(env *ExpressionEnv) (eval, error) {
	now := SystemTime()
	if tz := env.currentTimezone(); tz != nil {
		now = now.In(tz)
	}
	return newEvalRaw(sqltypes.Datetime, datetime.NewDateTimeFromStd(now).Format(call.prec), collationBinary), nil
}

func (call *builtinSysdate) typeof(_ *ExpressionEnv, _ []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Datetime, 0
}

func (call *builtinSysdate) compile(c *compiler) (ctype, error) {
	c.asm.Fn_Sysdate(call.prec)
	return ctype{Type: sqltypes.Datetime, Col: collationBinary}, nil
}

func (call *builtinSysdate) constant() bool {
	return false
}

func (call *builtinCurdate) eval(env *ExpressionEnv) (eval, error) {
	now := env.time(false)
	return newEvalRaw(sqltypes.Date, datetime.Date_YYYY_MM_DD.Format(now, 0), collationBinary), nil
}

func (call *builtinCurdate) typeof(_ *ExpressionEnv, _ []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Date, 0
}

func (*builtinCurdate) compile(c *compiler) (ctype, error) {
	c.asm.Fn_Curdate()
	return ctype{Type: sqltypes.Date, Col: collationBinary}, nil
}

func (call *builtinCurdate) constant() bool {
	return false
}

func (call *builtinUtcDate) eval(env *ExpressionEnv) (eval, error) {
	now := env.time(true)
	return newEvalRaw(sqltypes.Date, datetime.Date_YYYY_MM_DD.Format(now, 0), collationBinary), nil
}

func (call *builtinUtcDate) typeof(_ *ExpressionEnv, _ []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Date, 0
}

func (*builtinUtcDate) compile(c *compiler) (ctype, error) {
	c.asm.Fn_UtcDate()
	return ctype{Type: sqltypes.Date, Col: collationBinary}, nil
}

func (call *builtinUtcDate) constant() bool {
	return false
}

func (b *builtinDateFormat) eval(env *ExpressionEnv) (eval, error) {
	date, format, err := b.arg2(env)
	if err != nil {
		return nil, err
	}
	if date == nil || format == nil {
		return nil, nil
	}
	var t *evalTemporal
	switch e := date.(type) {
	case *evalTemporal:
		t = e.toDateTime(datetime.DefaultPrecision)
	default:
		t = evalToDateTime(date, datetime.DefaultPrecision)
		if t == nil || t.isZero() {
			return nil, nil
		}
	}

	f := evalToBinary(format)
	d, err := datetime.Format(f.string(), t.dt, t.prec)
	if err != nil {
		return nil, err
	}
	return newEvalText(d, defaultCoercionCollation(b.collate)), nil
}

func (b *builtinDateFormat) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.VarChar, flagNullable
}

func (call *builtinDateFormat) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip1 := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Datetime, sqltypes.Date:
	default:
		c.asm.Convert_xDT_nz(1, datetime.DefaultPrecision)
	}

	format, err := call.Arguments[1].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip2 := c.compileNullCheck1r(format)

	switch format.Type {
	case sqltypes.VarChar, sqltypes.VarBinary:
	default:
		c.asm.Convert_xb(1, sqltypes.VarBinary, 0, false)
	}

	col := defaultCoercionCollation(c.cfg.Collation)
	c.asm.Fn_DATE_FORMAT(col)
	c.asm.jumpDestination(skip1, skip2)
	return ctype{Type: sqltypes.VarChar, Col: col, Flag: arg.Flag | flagNullable}, nil
}

type builtinConvertTz struct {
	CallExpr
}

var _ Expr = (*builtinConvertTz)(nil)

func convertTz(dt datetime.DateTime, from, to *time.Location) (datetime.DateTime, bool) {
	buf := datetime.DateTime_YYYY_MM_DD_hh_mm_ss.Format(dt, datetime.DefaultPrecision)
	ts, err := time.ParseInLocation(time.DateTime, hack.String(buf), from)
	if err != nil {
		return datetime.DateTime{}, false
	}
	return datetime.NewDateTimeFromStd(ts.In(to)), true
}

func (call *builtinConvertTz) eval(env *ExpressionEnv) (eval, error) {
	n, err := call.Arguments[0].eval(env)
	if err != nil {
		return nil, err
	}
	from, err := call.Arguments[1].eval(env)
	if err != nil {
		return nil, err
	}
	to, err := call.Arguments[2].eval(env)
	if err != nil {
		return nil, err
	}

	if n == nil || from == nil || to == nil {
		return nil, nil
	}

	f := evalToBinary(from)
	t := evalToBinary(to)

	fromTz, err := datetime.ParseTimeZone(f.string())
	if err != nil {
		return nil, nil
	}

	toTz, err := datetime.ParseTimeZone(t.string())
	if err != nil {
		return nil, nil
	}

	dt := evalToDateTime(n, -1)
	if dt == nil || dt.isZero() {
		return nil, nil
	}

	out, ok := convertTz(dt.dt, fromTz, toTz)
	if !ok {
		return nil, nil
	}
	return newEvalDateTime(out, int(dt.prec)), nil
}

func (call *builtinConvertTz) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	_, f := call.Arguments[0].typeof(env, fields)
	return sqltypes.Datetime, f | flagNullable
}

func (call *builtinConvertTz) compile(c *compiler) (ctype, error) {
	n, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}
	from, err := call.Arguments[1].compile(c)
	if err != nil {
		return ctype{}, err
	}
	to, err := call.Arguments[2].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck3(n, from, to)

	switch {
	case from.isTextual():
	default:
		c.asm.Convert_xb(2, sqltypes.VarBinary, 0, false)
	}

	switch {
	case to.isTextual():
	default:
		c.asm.Convert_xb(1, sqltypes.VarBinary, 0, false)
	}

	switch n.Type {
	case sqltypes.Datetime, sqltypes.Date:
	default:
		c.asm.Convert_xDT_nz(3, -1)
	}

	c.asm.Fn_CONVERT_TZ()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Datetime, Col: collationBinary, Flag: n.Flag | flagNullable}, nil
}

func (b *builtinDate) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil {
		return nil, nil
	}
	return d, nil
}

func (b *builtinDate) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Date, flagNullable
}

func (call *builtinDate) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date:
	default:
		c.asm.Convert_xD(1)
	}

	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Date, Col: collationBinary, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinDayOfMonth) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Date.Day())), nil
}

func (b *builtinDayOfMonth) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinDayOfMonth) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD(1)
	}
	c.asm.Fn_DAYOFMONTH()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinDayOfWeek) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil || d.isZero() {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Date.ToStdTime(time.Local).Weekday() + 1)), nil
}

func (b *builtinDayOfWeek) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinDayOfWeek) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD_nz(1)
	}
	c.asm.Fn_DAYOFWEEK()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinDayOfYear) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil || d.isZero() {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Date.ToStdTime(time.Local).YearDay())), nil
}

func (b *builtinDayOfYear) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinDayOfYear) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD_nz(1)
	}
	c.asm.Fn_DAYOFYEAR()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

const maxUnixtime = 32536771200

func (b *builtinFromUnixtime) eval(env *ExpressionEnv) (eval, error) {
	ts, err := b.arg1(env)
	if err != nil {
		return nil, err
	}

	if ts == nil {
		return nil, nil
	}

	prec := 0
	var sec, frac int64

	switch ts := ts.(type) {
	case *evalInt64:
		sec = ts.i
	case *evalUint64:
		sec = int64(ts.u)
	case *evalFloat:
		sf, ff := math.Modf(ts.f)
		sec = int64(sf)
		frac = int64(ff * 1e9)
		prec = 6
	case *evalDecimal:
		sd, fd := ts.dec.QuoRem(decimal.New(1, 0), 0)
		sec, _ = sd.Int64()
		frac, _ = fd.Mul(decimal.New(1, 9)).Int64()
		prec = int(ts.length)
	case *evalTemporal:
		if ts.prec == 0 {
			sec = ts.toInt64()
		} else {
			dec := ts.toDecimal()
			sd, fd := dec.QuoRem(decimal.New(1, 0), 0)
			sec, _ = sd.Int64()
			frac, _ = fd.Mul(decimal.New(1, 9)).Int64()
			prec = int(ts.prec)
		}
	case *evalBytes:
		if ts.isHexOrBitLiteral() {
			u, _ := ts.toNumericHex()
			sec = int64(u.u)
		} else {
			f, _ := evalToFloat(ts)
			sf, ff := math.Modf(f.f)
			sec = int64(sf)
			frac = int64(ff * 1e9)
			prec = 6
		}
	default:
		f, _ := evalToFloat(ts)
		sf, ff := math.Modf(f.f)
		sec = int64(sf)
		frac = int64(ff * 1e9)
		prec = 6
	}

	if sec < 0 || sec >= maxUnixtime {
		return nil, nil
	}

	t := time.Unix(sec, frac)
	if tz := env.currentTimezone(); tz != nil {
		t = t.In(tz)
	}

	dt := newEvalDateTime(datetime.NewDateTimeFromStd(t), prec)

	if len(b.Arguments) == 1 {
		return dt, nil
	}

	format, err := b.Arguments[1].eval(env)
	if err != nil {
		return nil, err
	}
	f := evalToBinary(format)
	d, err := datetime.Format(f.string(), dt.dt, dt.prec)
	if err != nil {
		return nil, err
	}
	return newEvalText(d, defaultCoercionCollation(b.collate)), nil
}

func (b *builtinFromUnixtime) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	_, f := b.Arguments[0].typeof(env, fields)
	if len(b.Arguments) == 1 {
		return sqltypes.Datetime, f | flagNullable
	}
	return sqltypes.VarChar, f | flagNullable
}

func (call *builtinFromUnixtime) compile(c *compiler) (ctype, error) {
	arg, err := c.compile(call.Arguments[0])
	if err != nil {
		return ctype{}, err
	}
	skip1 := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Int64:
		c.asm.Fn_FROM_UNIXTIME_i()
	case sqltypes.Uint64:
		c.asm.Fn_FROM_UNIXTIME_u()
	case sqltypes.Float64:
		c.asm.Fn_FROM_UNIXTIME_f()
	case sqltypes.Decimal:
		c.asm.Fn_FROM_UNIXTIME_d()
	case sqltypes.Datetime, sqltypes.Date, sqltypes.Time:
		c.asm.Convert_Ti(1)
		c.asm.Fn_FROM_UNIXTIME_i()
	case sqltypes.VarChar, sqltypes.VarBinary:
		if arg.isHexOrBitLiteral() {
			c.asm.Convert_xu(1)
			c.asm.Fn_FROM_UNIXTIME_u()
		} else {
			c.asm.Convert_xf(1)
			c.asm.Fn_FROM_UNIXTIME_f()
		}
	default:
		c.asm.Convert_xf(1)
		c.asm.Fn_FROM_UNIXTIME_f()
	}

	if len(call.Arguments) == 1 {
		c.asm.jumpDestination(skip1)
		return ctype{Type: sqltypes.Datetime, Col: collationBinary, Flag: arg.Flag | flagNullable}, nil
	}

	format, err := call.Arguments[1].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip2 := c.compileNullCheck1r(format)

	switch format.Type {
	case sqltypes.VarChar, sqltypes.VarBinary:
	default:
		c.asm.Convert_xb(1, sqltypes.VarBinary, 0, false)
	}

	col := defaultCoercionCollation(c.cfg.Collation)
	c.asm.Fn_DATE_FORMAT(col)
	c.asm.jumpDestination(skip1, skip2)
	return ctype{Type: sqltypes.VarChar, Col: col, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinHour) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToTime(date, -1)
	if d == nil {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Time.Hour())), nil
}

func (b *builtinHour) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinHour) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime, sqltypes.Time:
	default:
		c.asm.Convert_xT(1, -1)
	}
	c.asm.Fn_HOUR()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func yearDayToTime(y, yd int64) time.Time {
	if y >= 0 && y < 100 {
		if y < 70 {
			y += 2000
		} else {
			y += 1900
		}
	}

	if y < 0 || y > 9999 || yd < 1 || yd > math.MaxInt32 {
		return time.Time{}
	}
	t := time.Date(int(y), time.January, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, int(yd-1))
	if t.Year() > 9999 {
		return time.Time{}
	}
	return t
}

func (b *builtinMakedate) eval(env *ExpressionEnv) (eval, error) {
	// For some reason, MySQL first evaluates the year day argument.
	yearDay, err := b.Arguments[1].eval(env)
	if err != nil {
		return nil, err
	}
	if yearDay == nil {
		return nil, nil
	}

	year, err := b.Arguments[0].eval(env)
	if err != nil {
		return nil, err
	}
	if year == nil {
		return nil, nil
	}

	y := evalToInt64(year).i
	yd := evalToInt64(yearDay).i

	t := yearDayToTime(y, yd)
	if t.IsZero() {
		return nil, nil
	}
	return newEvalDate(datetime.NewDateTimeFromStd(t).Date), nil
}

func (b *builtinMakedate) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	_, f1 := b.Arguments[0].typeof(env, fields)
	_, f2 := b.Arguments[1].typeof(env, fields)
	return sqltypes.Date, f1 | f2 | flagNullable
}

func (call *builtinMakedate) compile(c *compiler) (ctype, error) {
	// Similar here, we have to evaluate these in reverse order as well.
	yearDay, err := call.Arguments[1].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip1 := c.compileNullCheck1(yearDay)

	year, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip2 := c.compileNullCheck1r(year)

	switch yearDay.Type {
	case sqltypes.Int64:
	default:
		c.asm.Convert_xi(2)
	}

	switch year.Type {
	case sqltypes.Int64:
	default:
		c.asm.Convert_xi(1)
	}

	c.asm.Fn_MAKEDATE()
	c.asm.jumpDestination(skip1, skip2)
	return ctype{Type: sqltypes.Date, Col: collationBinary, Flag: year.Flag | yearDay.Flag | flagNullable}, nil
}

func clampHourMinute(h, m int64) (int64, int64, bool, bool) {
	var clamped bool
	if h > 838 || h < -838 {
		clamped = true
		if h > 0 {
			h = 838
		} else {
			h = -838
		}
		m = 59
	}
	neg := h < 0
	if neg {
		h = -h
	}
	return h, m, neg, clamped
}

func makeTime_i(h, m, s int64) (int64, bool) {
	if m < 0 || m > 59 || s < 0 || s > 59 {
		return 0, false
	}

	h, m, neg, clamped := clampHourMinute(h, m)
	if clamped {
		s = 59
	}

	v := h*10000 + m*100 + s
	if neg {
		v = -v
	}
	return v, true
}

func makeTime_d(h, m int64, s decimal.Decimal) (decimal.Decimal, bool) {
	if m < 0 || m > 59 || s.Sign() < 0 || s.Cmp(decimal.NewFromInt(60)) >= 0 {
		return decimal.Zero, false
	}

	h, m, neg, clamped := clampHourMinute(h, m)
	if clamped {
		s = decimal.NewFromInt(59)
	}

	dec := decimal.NewFromInt(h*10000 + m*100).Add(s)
	if neg {
		dec = dec.Neg()
	}
	return dec, true
}

func makeTime_f(h, m int64, s float64) (float64, bool) {
	if m < 0 || m > 59 || s < 0.0 || s >= 60.0 {
		return 0, false
	}

	h, m, neg, clamped := clampHourMinute(h, m)
	if clamped {
		s = 59.0
	}

	v := float64(h*10000+m*100) + s
	if neg {
		v = -v
	}
	return v, true
}

func (b *builtinMaketime) eval(env *ExpressionEnv) (eval, error) {
	hour, err := b.Arguments[0].eval(env)
	if err != nil {
		return nil, err
	}
	if hour == nil {
		return nil, nil
	}
	min, err := b.Arguments[1].eval(env)
	if err != nil {
		return nil, err
	}
	if min == nil {
		return nil, nil
	}
	sec, err := b.Arguments[2].eval(env)
	if err != nil {
		return nil, err
	}
	if sec == nil {
		return nil, nil
	}

	var h int64
	switch hour := hour.(type) {
	case *evalInt64:
		h = hour.i
	case *evalUint64:
		if hour.u > math.MaxInt64 {
			h = math.MaxInt64
		} else {
			h = int64(hour.u)
		}
	case *evalBytes:
		if hour.isHexOrBitLiteral() {
			hex, ok := hour.toNumericHex()
			if ok {
				if hex.u > math.MaxInt64 {
					h = math.MaxInt64
				} else {
					h = int64(hex.u)
				}
			}
		} else {
			h = evalToInt64(hour).i
		}
	default:
		h = evalToInt64(hour).i
	}

	m := evalToInt64(min).i
	s := evalToNumeric(sec, false)

	var ok bool
	var t datetime.Time
	var l int
	switch s := s.(type) {
	case *evalInt64:
		var v int64
		v, ok = makeTime_i(h, m, s.i)
		if !ok {
			return nil, nil
		}
		t, ok = datetime.ParseTimeInt64(v)
	case *evalUint64:
		var v int64
		v, ok = makeTime_i(h, m, int64(s.u))
		if !ok {
			return nil, nil
		}
		t, ok = datetime.ParseTimeInt64(v)
	case *evalDecimal:
		var v decimal.Decimal
		v, ok = makeTime_d(h, m, s.dec)
		if !ok {
			return nil, nil
		}
		t, l, ok = datetime.ParseTimeDecimal(v, s.length, -1)
	case *evalFloat:
		var v float64
		v, ok = makeTime_f(h, m, s.f)
		if !ok {
			return nil, nil
		}
		t, l, ok = datetime.ParseTimeFloat(v, -1)
	}
	if !ok {
		return nil, nil
	}

	return newEvalTime(t, l), nil
}

func (b *builtinMaketime) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	_, f1 := b.Arguments[0].typeof(env, fields)
	_, f2 := b.Arguments[1].typeof(env, fields)
	_, f3 := b.Arguments[2].typeof(env, fields)
	return sqltypes.Time, f1 | f2 | f3 | flagNullable
}

func (call *builtinMaketime) compile(c *compiler) (ctype, error) {
	hour, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}
	skip1 := c.compileNullCheck1(hour)

	min, err := call.Arguments[1].compile(c)
	if err != nil {
		return ctype{}, err
	}
	skip2 := c.compileNullCheck1r(min)

	sec, err := call.Arguments[2].compile(c)
	if err != nil {
		return ctype{}, err
	}
	skip3 := c.compileNullCheck3(hour, min, sec)

	switch hour.Type {
	case sqltypes.Int64:
	case sqltypes.Uint64:
		c.asm.Clamp_u(3, math.MaxInt64)
		c.asm.Convert_xi(3)
	case sqltypes.VarChar, sqltypes.VarBinary:
		if hour.isHexOrBitLiteral() {
			c.asm.Convert_xu(3)
			c.asm.Clamp_u(3, math.MaxInt64)
			c.asm.Convert_xi(3)
		} else {
			c.asm.Convert_xi(3)
		}
	default:
		c.asm.Convert_xi(3)
	}

	switch min.Type {
	case sqltypes.Int64:
	default:
		c.asm.Convert_xi(2)
	}

	switch sec.Type {
	case sqltypes.Int64:
		c.asm.Fn_MAKETIME_i()
	case sqltypes.Uint64:
		c.asm.Convert_ui(1)
		c.asm.Fn_MAKETIME_i()
	case sqltypes.Decimal:
		c.asm.Fn_MAKETIME_d()
	case sqltypes.Float64:
		c.asm.Fn_MAKETIME_f()
	case sqltypes.VarChar, sqltypes.VarBinary:
		if sec.isHexOrBitLiteral() {
			c.asm.Convert_xi(1)
			c.asm.Fn_MAKETIME_i()
		} else {
			c.asm.Convert_xf(1)
			c.asm.Fn_MAKETIME_f()
		}
	default:
		c.asm.Convert_xf(1)
		c.asm.Fn_MAKETIME_f()
	}

	c.asm.jumpDestination(skip1, skip2, skip3)
	return ctype{Type: sqltypes.Time, Col: collationBinary, Flag: hour.Flag | min.Flag | sec.Flag | flagNullable}, nil
}

func (b *builtinMicrosecond) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToTime(date, -1)
	if d == nil {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Time.Nanosecond() / 1000)), nil
}

func (b *builtinMicrosecond) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinMicrosecond) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime, sqltypes.Time:
	default:
		c.asm.Convert_xT(1, -1)
	}
	c.asm.Fn_MICROSECOND()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinMinute) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToTime(date, -1)
	if d == nil {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Time.Minute())), nil
}

func (b *builtinMinute) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinMinute) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime, sqltypes.Time:
	default:
		c.asm.Convert_xT(1, -1)
	}
	c.asm.Fn_MINUTE()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinMonth) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Date.Month())), nil
}

func (b *builtinMonth) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinMonth) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD(1)
	}
	c.asm.Fn_MONTH()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinMonthName) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil {
		return nil, nil
	}
	m := d.dt.Date.Month()
	if m < 1 || m > 12 {
		return nil, nil
	}

	return newEvalText(hack.StringBytes(time.Month(d.dt.Date.Month()).String()), defaultCoercionCollation(b.collate)), nil
}

func (b *builtinMonthName) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.VarChar, flagNullable
}

func (call *builtinMonthName) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD(1)
	}
	col := defaultCoercionCollation(call.collate)
	c.asm.Fn_MONTHNAME(col)
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.VarChar, Col: col, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinQuarter) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Date.Quarter())), nil
}

func (b *builtinQuarter) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinQuarter) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD(1)
	}
	c.asm.Fn_QUARTER()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinSecond) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToTime(date, -1)
	if d == nil {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Time.Second())), nil
}

func (b *builtinSecond) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinSecond) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime, sqltypes.Time:
	default:
		c.asm.Convert_xT(1, -1)
	}
	c.asm.Fn_SECOND()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinTime) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToTime(date, -1)
	if d == nil {
		return nil, nil
	}
	return d, nil
}

func (b *builtinTime) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Time, flagNullable
}

func (call *builtinTime) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Time:
	default:
		c.asm.Convert_xT(1, -1)
	}
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Time, Col: collationBinary, Flag: arg.Flag | flagNullable}, nil
}

func dateTimeUnixTimestamp(env *ExpressionEnv, date eval) evalNumeric {
	var dt *evalTemporal
	switch e := date.(type) {
	case *evalTemporal:
		dt = e.toDateTime(int(e.prec))
	default:
		dt = evalToDateTime(date, -1)
		if dt == nil || dt.isZero() {
			var prec int32
			switch d := date.(type) {
			case *evalInt64, *evalUint64:
				return newEvalInt64(0)
			case *evalDecimal:
				prec = d.length
			case *evalBytes:
				if d.isHexLiteral {
					return newEvalInt64(0)
				}
				prec = 6
			default:
				prec = 6
			}
			return newEvalDecimalWithPrec(decimal.Zero, prec)
		}
	}

	tz := env.currentTimezone()
	if tz == nil {
		tz = time.Local
	}

	ts := dt.dt.ToStdTime(tz)
	if dt.prec == 0 {
		return newEvalInt64(ts.Unix())
	}
	dec := decimal.New(ts.Unix(), 0)
	dec = dec.Add(decimal.New(int64(dt.dt.Time.Nanosecond()), -9))
	return newEvalDecimalWithPrec(dec, int32(dt.prec))
}

func (b *builtinUnixTimestamp) eval(env *ExpressionEnv) (eval, error) {
	if len(b.Arguments) == 0 {
		return newEvalInt64(env.now.Unix()), nil
	}

	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}

	if date == nil {
		return nil, nil
	}

	return dateTimeUnixTimestamp(env, date), nil
}

func (b *builtinUnixTimestamp) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	if len(b.Arguments) == 0 {
		return sqltypes.Int64, 0
	}
	_, f := b.Arguments[0].typeof(env, fields)
	return sqltypes.Int64, f | flagAmbiguousType
}

func (call *builtinUnixTimestamp) constant() bool {
	if len(call.Arguments) == 0 {
		return false
	}
	return call.Arguments[0].constant()
}

func (call *builtinUnixTimestamp) compile(c *compiler) (ctype, error) {
	if len(call.Arguments) == 0 {
		c.asm.Fn_UNIX_TIMESTAMP0()
		return ctype{Type: sqltypes.Int64, Col: collationNumeric}, nil
	}

	arg, err := c.compile(call.Arguments[0])
	if err != nil {
		return ctype{}, err
	}
	skip := c.compileNullCheck1(arg)

	c.asm.Fn_UNIX_TIMESTAMP1()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationBinary, Flag: arg.Flag | flagAmbiguousType}, nil
}

func (b *builtinWeek) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}

	d := evalToDate(date)
	if d == nil || d.isZero() {
		return nil, nil
	}

	mode := int64(0)
	if len(b.Arguments) == 2 {
		m, err := b.Arguments[1].eval(env)
		if err != nil {
			return nil, err
		}
		if m == nil {
			return nil, nil
		}
		mode = evalToInt64(m).i
	}

	week := d.dt.Date.Week(int(mode))
	return newEvalInt64(int64(week)), nil
}

func (b *builtinWeek) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinWeek) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip1 := c.compileNullCheck1(arg)
	var skip2 *jump

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD_nz(1)
	}

	if len(call.Arguments) == 1 {
		c.asm.Fn_WEEK0()
	} else {
		mode, err := call.Arguments[1].compile(c)
		if err != nil {
			return ctype{}, err
		}
		skip2 = c.compileNullCheck1r(mode)
		c.asm.Convert_xi(1)
		c.asm.Fn_WEEK()
	}

	c.asm.jumpDestination(skip1, skip2)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinWeekDay) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil || d.isZero() {
		return nil, nil
	}
	return newEvalInt64(int64(d.dt.Date.Weekday()+6) % 7), nil
}

func (b *builtinWeekDay) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinWeekDay) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD_nz(1)
	}

	c.asm.Fn_WEEKDAY()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinWeekOfYear) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil || d.isZero() {
		return nil, nil
	}

	_, week := d.dt.Date.ISOWeek()
	return newEvalInt64(int64(week)), nil
}

func (b *builtinWeekOfYear) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinWeekOfYear) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD_nz(1)
	}

	c.asm.Fn_WEEKOFYEAR()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinYear) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}
	d := evalToDate(date)
	if d == nil {
		return nil, nil
	}

	return newEvalInt64(int64(d.dt.Date.Year())), nil
}

func (b *builtinYear) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinYear) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip := c.compileNullCheck1(arg)

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD(1)
	}

	c.asm.Fn_YEAR()
	c.asm.jumpDestination(skip)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func (b *builtinYearWeek) eval(env *ExpressionEnv) (eval, error) {
	date, err := b.arg1(env)
	if err != nil {
		return nil, err
	}
	if date == nil {
		return nil, nil
	}

	d := evalToDate(date)
	if d == nil || d.isZero() {
		return nil, nil
	}

	mode := int64(0)
	if len(b.Arguments) == 2 {
		m, err := b.Arguments[1].eval(env)
		if err != nil {
			return nil, err
		}
		if m == nil {
			return nil, nil
		}
		mode = evalToInt64(m).i
	}

	week := d.dt.Date.YearWeek(int(mode))
	return newEvalInt64(int64(week)), nil
}

func (b *builtinYearWeek) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	return sqltypes.Int64, flagNullable
}

func (call *builtinYearWeek) compile(c *compiler) (ctype, error) {
	arg, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	skip1 := c.compileNullCheck1(arg)
	var skip2 *jump

	switch arg.Type {
	case sqltypes.Date, sqltypes.Datetime:
	default:
		c.asm.Convert_xD_nz(1)
	}

	if len(call.Arguments) == 1 {
		c.asm.Fn_YEARWEEK0()
	} else {
		mode, err := call.Arguments[1].compile(c)
		if err != nil {
			return ctype{}, err
		}
		skip2 = c.compileNullCheck1r(mode)
		c.asm.Convert_xi(1)
		c.asm.Fn_YEARWEEK()
	}

	c.asm.jumpDestination(skip1, skip2)
	return ctype{Type: sqltypes.Int64, Col: collationNumeric, Flag: arg.Flag | flagNullable}, nil
}

func evalToInterval(itv eval, unit datetime.IntervalType, negate bool) *datetime.Interval {
	switch itv := itv.(type) {
	case *evalBytes:
		return datetime.ParseInterval(itv.string(), unit, negate)
	case *evalFloat:
		return datetime.ParseIntervalFloat(itv.f, unit, negate)
	case *evalDecimal:
		return datetime.ParseIntervalDecimal(itv.dec, itv.length, unit, negate)
	default:
		return datetime.ParseIntervalInt64(evalToNumeric(itv, false).toInt64().i, unit, negate)
	}
}

func (call *builtinDateMath) eval(env *ExpressionEnv) (eval, error) {
	date, err := call.Arguments[0].eval(env)
	if err != nil || date == nil {
		return date, err
	}

	itv, err := call.Arguments[1].eval(env)
	if err != nil || itv == nil {
		return itv, err
	}

	interval := evalToInterval(itv, call.unit, call.sub)
	if interval == nil {
		return nil, nil
	}

	if tmp, ok := date.(*evalTemporal); ok {
		return tmp.addInterval(interval, collations.TypedCollation{}), nil
	}

	if tmp := evalToTemporal(date); tmp != nil {
		return tmp.addInterval(interval, defaultCoercionCollation(call.collate)), nil
	}

	return nil, nil
}

func (call *builtinDateMath) typeof(env *ExpressionEnv, fields []*querypb.Field) (sqltypes.Type, typeFlag) {
	tt, f := call.Arguments[0].typeof(env, fields)

	switch {
	case tt == sqltypes.Date && !call.unit.HasTimeParts():
		return sqltypes.Date, f | flagNullable
	case tt == sqltypes.Time && !call.unit.HasDateParts():
		return sqltypes.Time, f | flagNullable
	case tt == sqltypes.Datetime || tt == sqltypes.Timestamp || (tt == sqltypes.Date && call.unit.HasTimeParts()) || (tt == sqltypes.Time && call.unit.HasDateParts()):
		return sqltypes.Datetime, f | flagNullable
	default:
		return sqltypes.Char, f | flagNullable
	}
}

func (call *builtinDateMath) compile(c *compiler) (ctype, error) {
	date, err := call.Arguments[0].compile(c)
	if err != nil {
		return ctype{}, err
	}

	// TODO: constant propagation
	_, err = call.Arguments[1].compile(c)
	if err != nil {
		return ctype{}, err
	}

	var ret ctype
	ret.Flag = date.Flag | flagNullable
	ret.Col = collationBinary

	switch {
	case date.Type == sqltypes.Date && !call.unit.HasTimeParts():
		ret.Type = sqltypes.Date
		c.asm.Fn_DATEADD_D(call.unit, call.sub)
	case date.Type == sqltypes.Time && !call.unit.HasDateParts():
		ret.Type = sqltypes.Time
		c.asm.Fn_DATEADD_D(call.unit, call.sub)
	case date.Type == sqltypes.Datetime || date.Type == sqltypes.Timestamp || (date.Type == sqltypes.Date && call.unit.HasTimeParts()) || (date.Type == sqltypes.Time && call.unit.HasDateParts()):
		ret.Type = sqltypes.Datetime
		c.asm.Fn_DATEADD_D(call.unit, call.sub)
	default:
		ret.Type = sqltypes.VarChar
		ret.Col = defaultCoercionCollation(c.cfg.Collation)
		c.asm.Fn_DATEADD_s(call.unit, call.sub, ret.Col)
	}
	return ret, nil
}
