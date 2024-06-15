package model

import (
	"encoding/base64"
	"reflect"
	"strconv"
	"strings"

	"github.com/Yoshioka9709/yy-go-backend-template/util/errs"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs/codes"
	"github.com/Yoshioka9709/yy-go-backend-template/util/ptr"
	"golang.org/x/xerrors"
)

const cursorPrefix = "cursor:"

// Cursor カーソル
type Cursor string

// NewCursor カーソル文字列を生成
func NewCursor(pos int) Cursor {
	str := cursorPrefix + strconv.Itoa(pos)
	return Cursor(base64.StdEncoding.EncodeToString([]byte(str)))
}

// Pos カーソルから何番目かを取得
func (c Cursor) Pos() (int, error) {
	buf, err := base64.StdEncoding.DecodeString(string(c))
	if err != nil {
		return 0, err
	}
	if !strings.HasPrefix(string(buf), cursorPrefix) {
		return 0, xerrors.New("invalid cursor")
	}
	return strconv.Atoi(string(buf[len(cursorPrefix):]))
}

type PageInfo struct {
	EndCursor   *Cursor `json:"endCursor"`
	HasNextPage bool    `json:"hasNextPage"`
}

// Edge エッジインターフェイス
type Edge interface {
	Set(cursor Cursor, node any)
}

// PagerFunc ページャーの中身を返す関数
type PagerFunc func(limit, offset int) (any, error)

// Pager ページャー
type Pager struct {
	offset    int
	nodes     any
	nodeCount int
}

// NewForwardPager 前方ページャーを生成する
//
// `pagerFunc` にページャーの中身をスライスで返す関数を指定してください。
func NewForwardPager(dataPage *DataPage, pagerFunc PagerFunc) (*Pager, error) {
	p := &Pager{}
	// カーソルから何番目かを取得
	if dataPage.Offset != nil {
		p.offset = *dataPage.Offset
	}
	if dataPage.After != nil {
		after := Cursor(*dataPage.After)
		var err error
		p.offset, err = after.Pos()
		if err != nil {
			return nil, errs.Wrap(codes.InvalidParameter, err)
		}
	}

	if dataPage.First > 0 {
		var err error
		p.nodes, err = pagerFunc(dataPage.First, p.offset)
		if err != nil {
			return nil, err
		}
		rv := reflect.ValueOf(p.nodes)
		if rv.Kind() != reflect.Slice {
			return nil, xerrors.New("invalid node type")
		}
		p.nodeCount = rv.Len()
		if p.nodeCount > dataPage.First {
			return nil, xerrors.New("too many node count")
		}
	}

	return p, nil
}

// Edges 型を指定してエッジを返す
func (p *Pager) Edges(edgeType Edge) any {
	if edgeType == nil {
		return nil
	}

	rt := reflect.TypeOf(edgeType)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	// ノードが0件の場合は空のエッジを返す
	nodes := reflect.ValueOf(p.nodes)
	edges := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(rt)), p.nodeCount, p.nodeCount)
	for i := 0; i < p.nodeCount; i++ {
		edge := reflect.New(rt)
		cursor := NewCursor(p.offset + i + 1)
		node := nodes.Index(i)
		edge.Interface().(Edge).Set(cursor, node.Interface())
		edges.Index(i).Set(edge)
	}
	return edges.Interface()
}

// PageInfo ページング情報を返す
//
// 検索条件適用時の総件数を渡してください
func (p *Pager) PageInfo(totalCount int) *PageInfo {
	hasNextPage := p.offset+p.nodeCount < totalCount
	var endCursor *Cursor
	if p.nodeCount > 0 {
		endCursor = ptr.Ptr(NewCursor(p.offset + p.nodeCount))
	}

	return &PageInfo{
		EndCursor:   endCursor,
		HasNextPage: hasNextPage,
	}
}

// DataPage ページングで送信されるインプット
type DataPage struct {
	// 表示件数
	First int
	// 表示開始位置
	After *string
	// 表示開始位置(CMS)
	// NOTE: afterとoffsetの両方を指定した場合はafterが優先
	Offset *int
}
