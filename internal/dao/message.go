package dao

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/Confialink/wallet-messages/internal/model"

	"github.com/jinzhu/gorm"
)

const TypeIncoming = "incoming"
const TypeOutgoing = "outgoing"

type Message struct {
	db *gorm.DB
}

// Pagination is the abstract paginator
type Pagination struct {
	Limit  int
	Offset int
}

// NewRepository creates new repository
func NewMessage(db *gorm.DB) *Message {
	return &Message{db}
}

// FindByParams retrieve the list of messages
func (d *Message) FindByParams(params url.Values, userId string) ([]*model.Message, error) {
	var messages []*model.Message

	order := "last_message_created_at DESC" //p.Order

	if len(params.Get("sortField")) > 0 &&
		len(params.Get("sortDir")) > 0 {
		order = params.Get("sortField") + " " + params.Get("sortDir")
	}

	query := d.db.Debug()
	query = d.joinsChildren(query)
	query = d.selectMessages(query)
	query = d.filterNotDeleted(query, userId)
	query = query.Where("messages.parent_id IS NULL")
	query = d.filterByQueryParams(query, params)
	query = d.paginate(query, params)

	if err := query.
		Order(order).
		Find(&messages).
		Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// CountByParams retrieve count of messages by params
func (d *Message) CountByParams(params url.Values, userId string) (*int64, error) {
	var messages []*model.Message
	var count int64
	query := d.db

	query = d.joinsChildren(query)
	query = d.filterNotDeleted(query, userId)
	query = d.filterByQueryParams(query, params)

	if err := query.
		Find(&messages).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &count, nil
}

// FindByUserAndParams retrieve the list of messages
func (d *Message) FindByUserAndParams(userId string, params url.Values) ([]*model.Message, error) {
	var messages []*model.Message

	order := "last_message_created_at DESC" //p.Order

	if len(params.Get("sortField")) > 0 &&
		len(params.Get("sortDir")) > 0 {
		order = params.Get("sortField") + " " + params.Get("sortDir")
	}

	query := d.db.Debug()
	query = d.joinsChildren(query)
	query = d.selectMessages(query)
	query = d.filterNotDeleted(query, userId)
	query = d.filterByUser(query, userId, params)
	query = d.filterByType(query, userId, params)
	query = d.filterByParent(query, userId, params)
	query = d.filterUnread(query, userId, params)
	query = d.filterByQueryParams(query, params)
	query = d.paginate(query, params)

	if err := query.
		Order(order).
		Find(&messages).
		Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// FindAllUsers retrieve the list of users
func (d *Message) CountByUserAndParams(userId string, params url.Values) (*int64, error) {
	var messages []*model.Message
	var count int64
	query := d.db
	query = d.joinsChildren(query)
	query = d.filterNotDeleted(query, userId)
	query = d.filterByUser(query, userId, params)
	query = d.filterByType(query, userId, params)
	query = d.filterByParent(query, userId, params)
	query = d.filterByQueryParams(query, params)
	query = d.filterUnread(query, userId, params)

	if err := query.
		Find(&messages).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &count, nil
}

// FindByUserAndParams retrieve the list of messages
func (d *Message) FindUnassignedAndIncoming(params url.Values, userId string) ([]*model.Message, error) {
	var messages []*model.Message

	order := "last_message_created_at DESC" //p.Order

	if len(params.Get("sortField")) > 0 &&
		len(params.Get("sortDir")) > 0 {
		order = params.Get("sortField") + " " + params.Get("sortDir")
	}

	query := d.db

	query = d.joinsChildren(query)
	query = d.selectMessages(query)
	query = d.filterNotDeletedForAdmin(query, userId)
	query = d.filterByQueryParams(query, params)
	query = query.Where("messages.parent_id IS NULL AND messages.recipient_id IS NULL OR messages.parent_id IS NULL AND messages.is_recipient_incoming IS TRUE AND messages.recipient_id = ? OR messages.parent_id IS NULL AND messages.is_recipient_incoming IS NOT TRUE AND messages.sender_id = ?", userId, userId)
	query = d.filterUnreadForAdmin(query, userId, params)
	query = d.paginate(query, params)

	if err := query.
		Order(order).
		Find(&messages).
		Error; err != nil {
		return nil, err
	}

	return messages, nil
}

// FindAllUsers retrieve the list of users
func (d *Message) CountUnassignedAndIncoming(params url.Values, userId string) (*int64, error) {
	var messages []*model.Message
	var count int64

	query := d.db
	query = d.joinsChildren(query)
	query = d.filterNotDeletedForAdmin(query, userId)
	query = d.filterByQueryParams(query, params)
	query = query.Where("messages.parent_id IS NULL AND messages.recipient_id IS NULL OR messages.parent_id IS NULL AND messages.is_recipient_incoming IS TRUE AND messages.recipient_id = ? OR messages.parent_id IS NULL AND messages.is_recipient_incoming IS NOT TRUE AND messages.sender_id = ?", userId, userId)
	query = d.filterUnreadForAdmin(query, userId, params)

	if err := query.
		Find(&messages).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &count, nil
}

func (d *Message) FindByIDWithChildren(id uint64, userId string) (*model.Message, error) {
	var message model.Message
	message.ID = id

	query := d.db
	query = d.joinsChildren(query)
	query = d.filterNotDeleted(query, userId)
	query = d.preloadChildren(query, userId)
	query.First(&message)

	if err := query.
		Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (d *Message) FindByIDWithChildrenForAdmin(id uint64, userId string) (*model.Message, error) {
	var message model.Message
	message.ID = id

	query := d.db.Debug()
	query = d.joinsChildren(query)
	query = d.filterNotDeletedForAdmin(query, userId)
	query = d.preloadChildren(query, userId)

	if err := query.
		First(&message).
		Error; err != nil {
		return nil, err
	}

	return &message, nil
}

// FindByID find user by id
func (d *Message) FindByID(id uint64, userId string, filterNotDeleted bool) (*model.Message, error) {
	var message model.Message

	query := d.db
	query = d.joinsChildren(query)
	if filterNotDeleted {
		query = d.filterNotDeleted(query, userId)
	}
	query = query.Where("messages.id = ?", id)
	query.First(&message)

	if err := query.
		Error; err != nil {
		return nil, err
	}
	return &message, nil
}

// FindByID find an unread message by subject and sender
func (d *Message) FindUnreadBySubjectAndSender(recipientId string, senderId string, subject string) (*model.Message, error) {
	var message model.Message

	query := d.db
	query = query.Where("messages.recipient_id = ? AND messages.sender_id = ? AND messages.subject = ? AND messages.is_recipient_read = ?", recipientId, senderId, subject, 0)
	query.First(&message)

	if err := query.
		Error; err != nil {
		return nil, err
	}
	return &message, nil
}

// Create creates new user
func (d *Message) Create(message *model.Message) (*model.Message, error) {
	if err := d.db.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// BulkCreate creates many messages
func (d *Message) BulkCreate(messages []*model.Message) ([]*model.Message, error) {
	tx := d.db.Begin()
	for _, message := range messages {
		if err := tx.Create(message).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	return messages, nil
}

// Updates updates an existing user
func (d *Message) Updates(message *model.Message) (*model.Message, error) {
	if err := d.db.Model(message).Updates(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// Update updates an existing user
func (d *Message) Update(message *model.Message) (*model.Message, error) {
	if err := d.db.Model(message).Save(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// Delete delete an existing user
func (d *Message) Delete(message *model.Message) error {
	if err := d.db.Delete(message).Error; err != nil {
		return err
	}
	return nil
}

// CountUnreadByUser returns count of unread messages by user uid
func (d *Message) CountUnreadByUser(userId string) (*int64, error) {
	var count int64
	query := d.db.Model(&model.Message{})
	query = d.joinsChildren(query)
	query = d.filterNotDeleted(query, userId)
	if err := query.
		Select("COUNT(DISTINCT messages.id)").
		Where("messages.parent_id IS NULL AND messages.recipient_id = ? AND messages.is_recipient_read IS NOT TRUE", userId).
		Or("messages.parent_id IS NULL AND messages.sender_id = ? AND messages.is_sender_read IS NOT TRUE", userId).
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &count, nil
}

// CountUnreadByUser returns count of unread messages by user uid
func (d *Message) CountUnreadByAdmin(userId string) (*int64, error) {
	var count int64
	query := d.db.Model(&model.Message{})
	query = d.joinsChildren(query)
	query = d.filterNotDeleted(query, userId)
	if err := query.
		Select("COUNT(DISTINCT messages.id)").
		Where("messages.parent_id IS NULL AND messages.recipient_id = ? AND messages.is_recipient_read IS NOT TRUE", userId).
		Or("messages.parent_id IS NULL AND messages.sender_id = ? AND messages.is_sender_read IS NOT TRUE", userId).
		Or("messages.parent_id IS NULL and messages.recipient_id IS NULL AND messages.is_recipient_read IS NOT TRUE").
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	return &count, nil
}

func (d *Message) filterByQueryParams(query *gorm.DB, params url.Values) *gorm.DB {
	query = query.Group("messages.id")

	if params.Get("searchField") == "Subject" &&
		len(params.Get("searchQuery")) > 0 {
		query = query.Where("messages.subject LIKE ?", "%"+params.Get("searchQuery")+"%")
	}

	if params.Get("searchField") == "Message" &&
		len(params.Get("searchQuery")) > 0 {
		query = query.Where("messages.message LIKE ? OR child.message LIKE ?", "%"+params.Get("searchQuery")+"%", "%"+params.Get("searchQuery")+"%")
	}

	if params.Get("searchField") == "All" &&
		len(params.Get("searchQuery")) > 0 {
		var administratorTo = "administrator"
		var systemTo = "system"

		// TODO: these functions join "users.users" table which does not exist
		//query = d.joinsRecipient(query)
		//query = d.joinsSender(query)
		searchQuery := strings.ToLower(params.Get("searchQuery"))
		whereConditions := "messages.message LIKE ? OR messages.subject LIKE ? OR child.message LIKE ? "
		// TODO: see a comment above about joins
		//"OR recipient.first_name LIKE ? OR recipient.last_name LIKE ? " +
		//"OR sender.first_name LIKE ? OR sender.last_name LIKE ?"

		whereParams := []interface{}{
			"%" + searchQuery + "%",
			"%" + searchQuery + "%",
			"%" + searchQuery + "%",
			//"%" + searchQuery + "%",
			//"%" + searchQuery + "%",
			//"%" + searchQuery + "%",
			//"%" + searchQuery + "%",
		}
		if strings.Contains(administratorTo, searchQuery) {
			whereConditions += " OR messages.recipient_id is null"
		}
		if strings.Contains(systemTo, searchQuery) {
			whereConditions += " OR messages.sender_id = ''"
		}
		query = query.Where(whereConditions, whereParams...)
	}

	if len(params.Get("dateFrom")) > 0 {
		dateFrom := params.Get("dateFrom") + " 00:00:00"
		query = query.Where("messages.created_at > ? OR "+
			"child.created_at > ?", dateFrom, dateFrom)
	}

	if len(params.Get("dateTo")) > 0 {
		dateTo := params.Get("dateTo") + " 23:59:59"
		query = query.Where("messages.created_at < ? OR "+
			"child.created_at < ?", dateTo, dateTo)
	}

	return query
}

func (d *Message) filterUnread(query *gorm.DB, userId string, params url.Values) *gorm.DB {
	if len(params.Get("isUnread")) > 0 {
		query = query.
			Where("messages.parent_id IS NULL AND messages.recipient_id = ? AND "+
				"messages.is_recipient_read IS NOT TRUE OR "+
				"messages.parent_id IS NULL AND messages.sender_id = ? AND "+
				"messages.is_sender_read IS NOT TRUE", userId, userId)
	}
	return query
}

func (d *Message) filterUnreadForAdmin(query *gorm.DB, userId string, params url.Values) *gorm.DB {
	if len(params.Get("isUnread")) > 0 {
		query = query.
			Where("messages.parent_id IS NULL AND messages.recipient_id = ? AND "+
				"messages.is_recipient_read IS NOT TRUE OR "+
				"messages.parent_id IS NULL AND messages.sender_id = ? AND "+
				"messages.is_sender_read IS NOT TRUE OR "+
				"messages.parent_id IS NULL AND messages.recipient_id IS NULL AND messages.is_recipient_read IS NOT TRUE", userId, userId)
	}
	return query
}

// filterByUser filters data by recipient and sender
func (d *Message) filterByUser(query *gorm.DB, userId string, params url.Values) *gorm.DB {
	if params.Get("type") != TypeIncoming &&
		params.Get("type") != TypeOutgoing &&
		len(params.Get("parent")) == 0 {
		return query.
			Where("messages.parent_id IS NULL AND messages.recipient_id = ? "+
				"OR messages.parent_id IS NULL AND messages.sender_id = ?", userId, userId)
	}
	return query
}

// filterByType filters data by recipient or sender (incoming, outgoing)
func (d *Message) filterByType(query *gorm.DB, userId string, params url.Values) *gorm.DB {
	if params.Get("type") == TypeIncoming {
		return query.
			Where("messages.parent_id IS NULL AND messages.is_recipient_incoming IS TRUE AND messages.recipient_id = ? OR messages.parent_id IS NULL AND messages.is_recipient_incoming IS NOT TRUE AND messages.sender_id = ?", userId, userId)
	} else if params.Get("type") == TypeOutgoing {
		return query.
			Where("messages.parent_id IS NULL AND messages.is_recipient_incoming IS NOT TRUE AND messages.recipient_id = ? OR messages.parent_id IS NULL AND messages.is_recipient_incoming IS TRUE AND messages.sender_id = ?", userId, userId)
	}
	return query
}

// filterByParent filters data by parent field
func (d *Message) filterByParent(query *gorm.DB, userId string, params url.Values) *gorm.DB {
	if len(params.Get("parent")) > 0 {
		parent, _ := strconv.ParseInt(params.Get("parent"), 10, 32)
		return query.
			Where("messages.recipient_id = ? AND messages.parent_id = ?", userId, parent).
			Or("messages.sender_id = ? AND messages.parent_id = ?", userId, parent)
	}
	return query
}

func (d *Message) filterNotDeleted(query *gorm.DB, userId string) *gorm.DB {
	return query.
		Where("(messages.recipient_id = ? AND messages.deleted_for_recipient IS NOT TRUE) OR "+
			"(messages.sender_id = ? AND messages.deleted_for_sender IS NOT TRUE) OR "+
			"((child.recipient_id = ? AND child.deleted_for_recipient IS NOT TRUE) OR "+
			"(child.sender_id = ? AND child.deleted_for_sender IS NOT true))", userId, userId, userId, userId)
}

func (d *Message) filterNotDeletedForAdmin(query *gorm.DB, userId string) *gorm.DB {
	return query.
		Where("(messages.recipient_id = ? AND messages.deleted_for_recipient IS NOT TRUE) OR "+
			"(messages.sender_id = ? AND messages.deleted_for_sender IS NOT TRUE) OR "+
			"((child.recipient_id = ? AND child.deleted_for_recipient IS NOT TRUE) OR "+
			"(child.sender_id = ? AND child.deleted_for_sender IS NOT true)) OR "+
			"messages.recipient_id IS NULL OR child.recipient_id IS NULL", userId, userId, userId, userId)
}

func (d *Message) joinsChildren(query *gorm.DB) *gorm.DB {
	return query.Joins("LEFT JOIN messages child ON child.parent_id = messages.id")
}

func (d *Message) joinsRecipient(query *gorm.DB) *gorm.DB {
	return query.Joins("LEFT JOIN users.users recipient ON recipient.uid = messages.recipient_id")
}

func (d *Message) joinsSender(query *gorm.DB) *gorm.DB {
	return query.Joins("LEFT JOIN users.users sender ON sender.uid = messages.sender_id")
}

func (d *Message) preloadChildren(query *gorm.DB, userId string) *gorm.DB {
	return query.Preload("Children", "(recipient_id = ? AND "+
		"deleted_for_recipient IS NOT TRUE) OR "+
		"(sender_id = ? AND deleted_for_sender IS NOT TRUE)", userId, userId)
}

func (d *Message) selectMessages(query *gorm.DB) *gorm.DB {
	return query.Select("DISTINCT messages.id, messages.sender_id, messages.edited, messages.deleted_for_sender, " +
		"messages.deleted_for_recipient, messages.created_at, messages.updated_at, messages.is_sender_read, messages.is_recipient_read, " +
		"messages.is_recipient_incoming, messages.message, messages.subject, messages.recipient_id, messages.parent_id, messages.delete_after_read, " +
		"IF (MAX(child.created_at), MAX(child.created_at), messages.created_at) AS last_message_created_at").
		Group("messages.id")
}

// paginate check if query parameters limit, after and offset are set
// and applies them to query builder
// argument usedLimit will be set the same value as for limit
// it is used in order to determine whether there are more records exist
// limit always increments by 1
func (d *Message) paginate(query *gorm.DB, params url.Values) *gorm.DB {
	limit, err := strconv.ParseUint(params.Get("limit"), 10, 32)
	if nil == err {
		if limit > 100 {
			limit = 100
		}

	} else {
		limit = 10
	}

	query = query.Limit(uint(limit))
	offset, err := strconv.ParseUint(params.Get("offset"), 10, 32)
	if nil == err {
		query = query.Offset(uint(offset))
	}

	return query
}
