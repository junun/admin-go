import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, Steps, Radio,
 Select, Row, Col, Button, Popconfirm, Icon, Switch, DatePicker,
 Dropdown, Menu, message, Spin} from "antd";
import Editor from 'react-ace';
import {hasPermission, timeDatetimeTrans} from "@/utils/globalTools"
import {connect} from "dva";
import moment from 'moment';
import styles from './index.module.css';
import {cleanCommand, compareArray} from "@/utils/globalTools";
import Record from './Record';
import Info from './Info';


const dateFormat = 'YYYY-MM-DD';
const FormItem = Form.Item;

@connect(({ loading, schedule, host }) => {
    return {
      scheduleList: schedule.scheduleList,
      scheduleListLen: schedule.scheduleListLen,
      scheduleLoading: loading.effects['schedule/getSchedule'],
      hostList: host.hostList,
      scheduleInfo: schedule.scheduleInfo,
      scheduleInfoLoading: loading.effects['schedule/getScheduleInfo'],
    }
})


class SchedulePage extends React.Component {
  constructor(props) {
    super(props);
    this.isFirstRender = true;
  }

  state = {
    visible: false,
    infoVisible: false,
    recordVisible: false,
    editCacheData: {},
    step: 0,
    loading: false,
    specShow: "1",
    Command: '',
    info: {},
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({
      type: 'schedule/getSchedule',
      payload: {
        page: 1,
        pageSize: 10, 
      }
    });
    dispatch({ type: 'host/getHost',
      payload: {
        page: 1,
        pageSize: 9999,
        Status: 1,
        Source: "schedule",
      }
    });
  }

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'schedule/getSchedule',
      payload: {
        page: page,
        pageSize: 10, 
      }
    });
  };

  showTypeAddModal = () => {
    this.setState({ 
      editCacheData: {},
      Command: '',
      visible: true 
    });
  };

  handleCancel = () => {
    this.setState({
      visible: false,
      step: 0,
      Command: '',
      editCacheData: {},
      specShow: "1",
    });
    this.isFirstRender = true;
  };

  handleChange = (e) => {
    this.setState({
      specShow: e.target.value
    });
  };

  openNotification = () => {
    notification.open({
      message: 'Notification Title',
      description:
        'This is the content of the notification. This is the content of the notification. This is the content of the notification.',
      onClick: () => {
        console.log('Notification Clicked!');
      },
    });
  };

  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        values.Command = this.state.Command;
        const obj = this.state.editCacheData;
        values.TriggerType  =  parseInt(values.TriggerType)
        values.IsMore       =  values.IsMore ? 1 : 0

        if (Object.keys(obj).length) {
          values.id = obj.id;
          values.HostIds = values.HostIds.join(",")
          console.log(values)
          dispatch({
            type: 'schedule/scheduleEdit',
            payload: values,
          });
          this.setState({ 
            visible: false,
            step: 0,
            Command: '',
          });
          //判断情况太多，前端不做是否有内容更改检查！
          // var arr = JSON.parse("[" + obj.HostIds + "]");
          // if (compareArray(arr, values.HostIds)) {
          //   if (values.TriggerType == 2 ) {
          //     if (
          //       obj.Name        === values.Name &&
          //       obj.Spec        === values.Spec._i &&
          //       obj.IsMore      === values.IsMore &&
          //       obj.HostIds     === values.HostIds.join(",") &&
          //       obj.Command     === values.Command &&
          //       obj.TriggerType === values.TriggerType &&
          //       obj.Desc        === values.Desc 
          //     ) {
          //        message.warning('没有内容修改， 请检查。');
          //     } else {
          //       values.id = obj.id;
          //       values.HostIds = values.HostIds.join(",")
          //       console.log(values)
          //       dispatch({
          //         type: 'schedule/scheduleEdit',
          //         payload: values,
          //       });
          //       this.setState({ 
          //         visible: false,
          //         step: 0,
          //         Command: '',
          //       });
          //     }
          //   } else {
          //     values.id = obj.id;
          //     values.HostIds = values.HostIds.join(",")
          //     console.log(values)
          //     dispatch({
          //       type: 'schedule/scheduleEdit',
          //       payload: values,
          //     });
          //     this.setState({ 
          //       visible: false,
          //       step: 0,
          //       Command: '',
          //     });
          //   }
          // } else {
          //   values.id = obj.id;
          //   values.HostIds = values.HostIds.join(",")
          //   console.log(values)
          //   dispatch({
          //     type: 'schedule/scheduleEdit',
          //     payload: values,
          //   });
          //   this.setState({ 
          //     visible: false,
          //     step: 0,
          //     Command: '',
          //   });
          // }
        } else {
          values.HostIds = values.HostIds.join(",");
          dispatch({
            type: 'schedule/scheduleAdd',
            payload: values,
          });
          this.setState({ 
            visible: false,
            step: 0,
            Command: '',
          });
        }
        // 重置 `visible` 属性为 false 以关闭对话框
      }
    });
  };

  //显示编辑界面
  handleEdit = (values) => {
    values.title =  '编辑-' + values.Name;
    var specShowCheck = "1"
    if (values.TriggerType == 2) {
      specShowCheck = "2"
    }
    if (values.StartTime == "0001-01-01T00:00:00Z") {
      values.StartTime = null
    }

    if (values.EndTime == "0001-01-01T00:00:00Z") {
      values.EndTime = null
    }

    this.setState({ 
      visible: true ,
      editCacheData: values,
      Command: values.Command,
      specShow: specShowCheck,
    });
  };

  // 删除一条记录
  deleteRecord = (values) => {
    console.log(values);
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'schedule/scheduleDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  handCommand = (v) => {
    var cmd = cleanCommand(v);
    this.setState({ 
      Command: cmd,
    });
  };

  showRecord = (v) => {
    this.setState({ 
      info: v,
      recordVisible: true,
    });
  };

  showInfo = (type, id) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'schedule/getScheduleInfo',
      payload: {
        id: id,
        Type: type,
      }
    });
    this.setState({ 
      infoVisible: true,
    });
  };

  handleCancelInfo = () => {
    this.setState({
      infoVisible: false,
      recordVisible: false,
      info: {},
      infoType:1,
      infoId:0,
    });
  };

  handleActive = (text) => {
    Modal.confirm({
      title: '删除确认',
      content: `确定要${text.Active == 1 ? '禁用' : '激活'}任务【${text.Name}】?`,
      onOk: () => {
        const { dispatch } = this.props;
        var sts = 1
        if (text.Active == 1) {
          sts = 0
        }

        dispatch({
          type: 'schedule/changeScheduleActive',
          payload: {id: text.id, Active: sts, Name: text.Name},
        });
      }
    })
  };

  verifyButtonStatus = () => {
    const data = this.props.form.getFieldsValue();
    let   b1 = data.Name && this.state.Command;
    const b2 = data.HostIds != '';
    const b3 = data.Spec && data.TriggerType;
    if (!b1 && this.isFirstRender && Object.keys(this.state.editCacheData).length) {
      this.isFirstRender = false;
      b1 = true
    }
    return [b1, b2, b3];
  };

  verifyJobStatus = (Active, TriggerType, StartTime, EndTime, Spec) => {
    if (Active==0) {
      return "Job未激活或已暂停"
    }

    if (TriggerType==1) {
      if (moment().isAfter(Spec)) {
        return "Job已执行"
      } else {
        return "Job设定运行时间未到"
      }
    }

    if (StartTime != "0001-01-01T00:00:00Z") {
      if (!moment().isAfter(StartTime)) {
        return "Job未到开始运行时间"
      } 
    }

    if (EndTime != "0001-01-01T00:00:00Z") {
      if (moment().isAfter(EndTime)) {
        return "Job已过结束运行时间"
      } 
    }

    return "Job正常调度中"
  }

  columns = [
  {
    title: '任务名',
    dataIndex: 'Name',
  }, {
    title: '状态',
    dataIndex: 'Active',
    ellipsis: true,
    'render': (Active, TriggerType, StartTime, EndTime, Spec) => this.verifyJobStatus(Active, TriggerType, StartTime, EndTime, Spec),
  }, {
    title: '触发器',
    dataIndex: 'Spec',
  }, {
    title: '描述信息',
    dataIndex: 'Desc',
    ellipsis: true
  }, {
    title: '操作',
    width: 200,
    render: (text, record) => (
      <span>
        {hasPermission('schedule-job-info') && <a onClick={()=>{this.showInfo(1, record.id)}}><Icon type="message"/>详情</a>}
        <Divider type="vertical" />
          {hasPermission('schedule-job-edit') && <a onClick={()=>{this.handleEdit(record)}}><Icon type="edit"/>编辑</a>}
        <Divider type="vertical" />
        <Dropdown overlay={() => this.moreMenus(record)} trigger={['click']}>
          <a>
            更多 <Icon type="down"/>
          </a>
        </Dropdown>
      </span>
    ),
  }];

  moreMenus = (info) => (
    <Menu>
      <Menu.Item>
        {hasPermission('schedule-job-edit') && <a onClick={() => this.handleActive(info)}>{info.Active == 1 ? '禁用任务' : '激活任务'}</a>}
      </Menu.Item>
      <Menu.Item>
        <a onClick={() => this.showRecord(info)}>历史记录</a>
      </Menu.Item>
      <Menu.Divider/>
      <Menu.Item>
        <Popconfirm title="你确定要删除吗?"  onConfirm={()=>{this.deleteRecord(info.id)}}>
          {hasPermission('schedule-job-del') && <a title="删除" >删除</a>}
        </Popconfirm>
      </Menu.Item>
    </Menu>
  );
  
  render() {
    const {step, loading, specShow, visible, infoVisible,
    recordVisible, editCacheData, Command} = this.state;
    const {scheduleList, scheduleListLen, scheduleLoading, hostList,
     scheduleInfo, scheduleInfoLoading, form: { getFieldDecorator } } = this.props;
    const [b1, b2, b3] = this.verifyButtonStatus();
    const addScheduleRple = <Button type="primary" onClick={this.showTypeAddModal} >新增任务</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('schedule-info-add') && <Col span={10}>{addScheduleRple}</Col>}
      </Row>; 

    return (
      <div>
        <Modal
          visible={visible}
          width={800}
          destroyOnClose= "true"
          maskClosable={false}
          title={editCacheData.title || "新增任务"}
          onCancel={this.handleCancel}
          footer={null}>
          <Steps current={step} className={styles.steps}>
            <Steps.Step key={0} title="创建任务"/>
            <Steps.Step key={1} title="选择执行对象"/>
            <Steps.Step key={2} title="设置触发器"/>
          </Steps>
          <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
            <div style={{display: step === 0 ? 'block' : 'none'}}>
              <FormItem label="任务名称">
                {getFieldDecorator('Name', {
                  initialValue: editCacheData.Name || '',
                  rules: [{ required: true }],
                })(
                  <Input placeholder="请输入任务名称"/>
                )}
              </FormItem>
              <FormItem required label="任务内容">
                  <Editor
                    theme="tomorrow"
                    enableLiveAutocompletion={true}
                    enableBasicAutocompletion={true}
                    enableSnippets={true}
                    mode="sh"
                    value={ Command || ""}
                    onChange={v => this.handCommand(v)}
                    height="150px"/>
              </FormItem>
              <FormItem label="多实例支持">
                {getFieldDecorator('IsMore', {
                  initialValue: editCacheData.IsMore && true || false,
                  valuePropName: "checked",
                  rules: [{ required: false }],
                })(
                  <Switch/>
                )}
              </FormItem>
              <FormItem label="备注信息">
                {getFieldDecorator('Desc', {
                  initialValue: editCacheData.Desc || '',
                  rules: [{ required: false }],
                })(
                  <Input.TextArea />
                )}
              </FormItem>
            </div>
            <div style={{display: step === 1 ? 'block' : 'none'}}>
              <FormItem required label="执行目标主机">
                {getFieldDecorator('HostIds', {
                  initialValue: editCacheData.HostIds && JSON.parse("[" + editCacheData.HostIds + "]") || [],
                  rules: [{ required: false }],
                })(
                  <Select
                    mode="multiple"
                    style={{ width: '100%' }}
                    filterOption={(input, option) => option.props.children[0].toLowerCase().indexOf(input.toLowerCase()) >= 0}
                  >
                    {hostList.map(item => (
                      <Select.Option key={item.id} value={item.id}>
                        {item.Name}({item['Addres']}:{item['Port']})
                      </Select.Option>
                    ))}
                  </Select>
                )}
              </FormItem>
            </div>
            <div style={{display: step === 2 ? 'block' : 'none'}}>
              <FormItem label="trigger type">
                {getFieldDecorator('TriggerType', {
                  initialValue: editCacheData.TriggerType && editCacheData.TriggerType.toString() || "1",
                  rules: [{ required: true }],
                })(
                  <Radio.Group size="small" onChange={this.handleChange}>
                    <Radio.Button value="1">一次性</Radio.Button>
                    <Radio.Button value="2">UNIX Cron</Radio.Button>
                  </Radio.Group>
                )}
              </FormItem>
              { specShow == 1 && 
                <FormItem required label="执行时间" extra="仅在指定时间运行一次。">
                  {getFieldDecorator('Spec', {
                    initialValue: editCacheData.Spec && moment(editCacheData.Spec, dateFormat) || null,
                    rules: [{ required: false }],
                  })(
                    <DatePicker
                      showTime
                      disabledDate={v => v && v.format(dateFormat) < moment().format(dateFormat)}
                      style={{width: 150}}
                      placeholder="请选择执行时间"
                      onOk={() => false}
                    />
                  )}
                </FormItem>
              }
              { specShow == 2 && 
                <FormItem required label="执行规则" help="兼容Cron风格，可参考官方例子">
                  {getFieldDecorator('Spec', {
                    initialValue: editCacheData.Spec || '',
                    rules: [{ required: false }],
                  })(
                    <Input placeholder="例如每天分钟执行：*/1 * * * ?"/>
                  )}
                </FormItem>
              }
              { specShow == 2 && 
                <FormItem label="生效时间" help="定义的执行规则在到达该时间后生效">
                  {getFieldDecorator('StartTime', {
                    initialValue: editCacheData.StartTime && editCacheData.StartTime != "0001-01-01T00:00:00Z" && moment(editCacheData.StartTime, dateFormat) || null,
                    rules: [{ required: false }],
                  })(
                    <DatePicker
                      showTime
                      style={{width: '100%'}}
                      placeholder="可选输入"
                    />
                  )}
                </FormItem>
              }
              { specShow == 2 && 
                <FormItem label="结束时间" help="执行规则在到达该时间后不再执行">
                  {getFieldDecorator('EndTime', {
                    initialValue: editCacheData.EndTime && editCacheData.EndTime != "0001-01-01T00:00:00Z" && moment(editCacheData.EndTime, dateFormat) || null,
                    rules: [{ required: false }],
                  })(
                    <DatePicker
                      showTime
                      style={{width: '100%'}}
                      placeholder="可选输入"
                    />
                  )}
                </FormItem>
              }
            </div>
            <Form.Item wrapperCol={{span: 14, offset: 6}}>
              {step === 2 &&
              <Button disabled={!b3} type="primary" onClick={this.handleOk} loading={loading}>提交</Button>}
              {step === 0 &&
              <Button disabled={!b1} type="primary" onClick={() => this.setState({step: step + 1})}>下一步</Button>}
              {step === 1 &&
              <Button disabled={!b2} type="primary" onClick={() => this.setState({step: step + 1})}>下一步</Button>}
              {step !== 0 &&
              <Button style={{marginLeft: 20}} onClick={() => this.setState({step: step - 1})}>上一步</Button>}
            </Form.Item>
          </Form>
        </Modal>

        {
          recordVisible && 
          <Record 
            info={this.state.info}
            onCancel={this.handleCancelInfo}
            showInfo={this.showInfo}
            hostList={hostList}
          />
        }

        {
          infoVisible && 
          <Info 
            onCancel={this.handleCancelInfo}
            scheduleInfo={scheduleInfo}
            scheduleInfoLoading={scheduleInfoLoading}
            hostList={hostList}
          />
        }

        <Card title="" extra={extra}>
          <Table  
          pagination={{
            showQuickJumper: true,
            total: scheduleListLen,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={scheduleList} loading={scheduleLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(SchedulePage);
